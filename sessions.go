package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var (
	errInvalidPassword = errors.New("The email and password you entered did not match our records. Please check them and try again.")
	errLockedAccount   = errors.New("This account is locked. Please check your email for a message containing a link to unlock your account.")
)

func isLoggedIn(w http.ResponseWriter, r *http.Request) (isLoggedIn bool, err error) {
	userID, ok := cookieIntValue("user-id", r)
	if !ok {
		return
	}

	token, ok := cookieStringValue("auth-token", r)
	if !ok {
		return
	}

	err = db.QueryRow("SELECT user_id FROM user_sessions WHERE user_id = $1 AND token = $2", userID, token).Scan(&userID)

	if err == sql.ErrNoRows {
		logOutUser(w, r)
		err = nil
	} else if err != nil {
		return
	} else {
		isLoggedIn = true
	}
	return
}

func currentUserID(r *http.Request) (userID int, ok bool) {
	return cookieIntValue("user-id", r)
}

func cookieIntValue(cookieName string, r *http.Request) (value int, ok bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	value, err = strconv.Atoi(cookie.Value)
	if err != nil {
		return
	}
	ok = true
	return
}

func cookieStringValue(cookieName string, r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return
	}
	value = cookie.Value
	ok = true
	return
}

func logInAjax(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	err := logInUser(r.PostFormValue("email"), r.PostFormValue("password"), w)
	if err == errInvalidPassword || err == errLockedAccount {
		json.NewEncoder(w).Encode(ErrorJSON{
			Messages: []string{err.Error()},
			Fields:   []string{"email", "password"},
		})
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(ErrorJSON{
			Errors:   []string{"An error occurred while trying to log in."},
			Messages: []string{adminNotifiedMessage},
			Debug:    []string{err.Error()},
		})
		return
	}
	json.NewEncoder(w).Encode(SuccessJSON{Success: true})
}

func logInUser(email string, password string, w http.ResponseWriter) (err error) {
	userID, err := validEmailAndPasswordUserID(email, password)
	if err == errInvalidPassword || err == errLockedAccount {
		return
	}

	deleteExpiredSessions()

	token := generateStringToken(10)
	_, err = db.Query("INSERT INTO user_sessions (user_id, token) VALUES ($1, $2)", userID, token)

	if err != nil {
		return
	}

	newCookie := http.Cookie{Name: "user-id", Value: strconv.Itoa(userID), Path: "/"}
	http.SetCookie(w, &newCookie)
	newCookie = http.Cookie{Name: "auth-token", Value: string(token), Path: "/"}
	http.SetCookie(w, &newCookie)

	return
}

func validEmailAndPasswordUserID(email string, password string) (validUserID int, err error) {
	emailHash, err := scryptHash(email, []byte(siteData.EmailSalt))
	if err != nil {
		return
	}

	var userID int
	var passwordHash, passwordSalt []byte
	err = db.QueryRow("SELECT id, password_hash, password_salt FROM users WHERE email_hash = $1", emailHash).Scan(&userID, &passwordHash, &passwordSalt)
	if err == sql.ErrNoRows {
		err = errInvalidPassword
	}
	if err != nil {
		return
	}

	correctPasswordHash, err := scryptHash(password, passwordSalt)
	if err != nil {
		return
	}

	if string(correctPasswordHash) != string(passwordHash) {
		err = errInvalidPassword
		return
	}

	validUserID = userID

	return
}

func redirectToHomeIfLoggedIn(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	loggedIn, err := isLoggedIn(w, r)
	if err != nil {
		serve500(w)
		return
	}
	if loggedIn {
		http.Redirect(w, r, "/", 302)
		return
	}
	serveStaticFilesOr404(w, r)
}

func logOut(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	logOutUser(w, r)
	http.Redirect(w, r, "/login/", 302)
}

func logOutAjax(w http.ResponseWriter, r *http.Request, requestParameters httprouter.Params) {
	logOutUser(w, r)
	json.NewEncoder(w).Encode(SuccessJSON{Success: true})
}

func logOutUser(w http.ResponseWriter, r *http.Request) {
	userIDCookie, userIDErr := r.Cookie("user-id")
	tokenCookie, tokenErr := r.Cookie("auth-token")
	if userIDErr == nil && tokenErr == nil {
		_, err := db.Query("DELETE FROM user_sessions WHERE user_id = $1 AND token = $2", userIDCookie.Value, tokenCookie.Value)
		if err != nil {
			notifyAdmin(err.Error())
		}
	}
	deleteCookie := http.Cookie{Name: "user-id", MaxAge: -1}
	http.SetCookie(w, &deleteCookie)
	deleteCookie = http.Cookie{Name: "auth-token", MaxAge: -1}
	http.SetCookie(w, &deleteCookie)
}

func serveLoginPage(w http.ResponseWriter) {
	template, err := ioutil.ReadFile(webRoot + "/login/index.html")
	if err != nil {
		serve500(w)
	}
	fmt.Fprint(w, string(template))
}

func deleteExpiredSessions() {
	_, err := db.Query("DELETE FROM user_sessions WHERE created_at < NOW() - INTERVAL '" + strconv.Itoa(sessionDurationInDays) + " days'")
	if err != nil {
		notifyAdmin(err.Error())
	}
}

func createUser(email string, password string) (err error) {
	emailHash, err := scryptHash(email, []byte(siteData.EmailSalt))

	exists, err := emailIsRegistered(emailHash)
	if exists {
		err = errors.New("An account has already been registered with this email address.")
		return
	}

	passwordHash, passwordSalt, err := scryptHashAndSalt(password)

	_, err = db.Query("INSERT INTO users (email_hash, password_hash, password_salt) VALUES ($1, $2, $3)", emailHash, passwordHash, passwordSalt)
	if err != nil {
		return
	}

	err = sendEmail(email, "Priori Sign Up Success!!", "We did it!! Woohoo!!")

	return
}

func emailIsRegistered(emailHash []byte) (exists bool, err error) {
	var id uint64
	err = db.QueryRow("SELECT id FROM users WHERE email_hash = $1", emailHash).Scan(&id)
	if err == sql.ErrNoRows {
		err = nil
	} else if err != nil {
		return
	} else {
		exists = true
	}
	return
}
