package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type SiteData struct {
	LiveOrDev             string            `json:"live-or-dev"`
	NoReplyAddressName    string            `json:"no-reply-address-name"`
	NoReplyAddress        string            `json:"no-reply-address"`
	NoReplyPassword       string            `json:"no-reply-password"`
	Host                  string            `json:"no-reply-host"`
	Port                  string            `json:"no-reply-port"`
	ReplyAddress          string            `json:"reply-address"`
	AdminEmail            string            `json:"admin-email"`
	DatabaseName          string            `json:"database-name"`
	EmailSalt             string            `json:"email-salt"`
	URLPermanentRedirects map[string]string `json:"url-permanent-redirects"`
}

var (
	webRoot               = "awestruct/_site"
	siteData              = SiteData{}
	db                    *sql.DB
	sessionDurationInDays = 4
	adminNotifiedMessage  = "An admin has been notified. Please try again. If the issue persists, please contact us to let us know."
)

type apiResponse struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Fields   []string `json:"fields,omitempty"`
	Debug    []string `json:"debug,omitempty"`
}

func main() {
	loadSiteData()
	db = connectToDatabase()
	router := httprouter.New()
	// Allows requests to pass through to NotFound if one method
	// is there and the other is not
	router.HandleMethodNotAllowed = false

	router.GET("/", authorize(serveStaticFilesOr404Handler))

	router.GET("/logout/", logOut)
	router.GET("/login/", redirectToHomeIfLoggedIn)

	router.POST("/api-noauth/", apiNoauth)
	router.POST("/api/", apiAuthorize(api))

	router.NotFound = http.HandlerFunc(serveStaticFilesOr404)
	log.Fatal(http.ListenAndServe(":9000", router))
}

func apiNoauth(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	allowedFunctions := map[string]func(http.ResponseWriter, *http.Request) interface{}{
		"login":  apiLogIn,
		"logout": apiLogOut,
	}
	apiGeneral(w, r, allowedFunctions)
}

func api(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	allowedFunctions := map[string]func(http.ResponseWriter, *http.Request) interface{}{
		"notify-admin": apiNotifyAdmin,
		"tasks":        apiTasks,
		"new-task":     apiNewTask,
		"delete-task":  apiDeleteTask,
	}
	apiGeneral(w, r, allowedFunctions)
}

func apiGeneral(w http.ResponseWriter, r *http.Request, allowedFunctions map[string]func(http.ResponseWriter, *http.Request) interface{}) {
	if r.PostFormValue("action") == "" {
		serve404(w)
		return
	}
	function, ok := allowedFunctions[r.PostFormValue("action")]
	if !ok {
		json.NewEncoder(w).Encode(apiResponse{
			Errors: []string{"The requested action could not be found in our api."},
		})
		return
	}
	switch response := function(w, r).(type) {
	case apiResponse, []Task:
		json.NewEncoder(w).Encode(response)
	case string:
		fmt.Fprint(w, response)
	case int:
		responseString := strconv.Itoa(response)
		fmt.Fprint(w, responseString)
	default:
		err := errors.New("Api response type is of an unknown type.")
		json.NewEncoder(w).Encode(notifyAdminResponse("An error occurred while accessing the api.", err))
	}
}

func serveStaticFilesOr404Handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serveStaticFilesOr404(w, r)
}

func debug(things ...interface{}) {
	if siteData.LiveOrDev == "dev" {
		fmt.Println("====================")
		for _, thing := range things {
			fmt.Printf("%+v\n", thing)
		}
		fmt.Println("^^^^^^^^^^^^^^^^^^^^")
	}
}

func debugType(thing interface{}) {
	fmt.Printf("%T\n", thing)
}
