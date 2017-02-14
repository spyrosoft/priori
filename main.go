package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type SiteData struct {
	LiveOrDev             string            `json:"live-or-dev"`
	URLPermanentRedirects map[string]string `json:"url-permanent-redirects"`
	NoReplyAddressName    string            `json:"no-reply-address-name"`
	NoReplyAddress        string            `json:"no-reply-address"`
	NoReplyPassword       string            `json:"no-reply-password"`
	Host                  string            `json:"no-reply-host"`
	Port                  string            `json:"no-reply-port"`
	ReplyAddress          string            `json:"reply-address"`
	AdminEmail            string            `json:"admin-email"`
	DatabaseName          string            `json:"database-name"`
	EmailSalt             string            `json:"email-salt"`
}

var (
	webRoot               = "awestruct/_site"
	siteData              = SiteData{}
	siteDataLoaded        = false
	db                    *sql.DB
	sessionDurationInDays = 4
	adminNotifiedMessage  = "An admin has been notified. Please try again. If the issue persists, please contact us to let us know."
)

type ErrorJSON struct {
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
	Fields   []string `json:"fields"`
	Debug    []string `json:"debug"`
}

type SuccessJSON struct {
	Success bool `json:"success"`
}

func main() {
	loadSiteData()
	db = connectToDatabase()
	router := httprouter.New()

	// Allows requests to pass through to NotFound if one method
	// is there and the other is not
	router.HandleMethodNotAllowed = false

	router.GET("/", authorize(serveStaticFilesOr404Handler))

	router.POST("/logout/", logOutAjax)
	router.GET("/logout/", logOut)

	router.POST("/login/", logInAjax)
	router.GET("/login/", redirectToHomeIfLoggedIn)

	router.POST("/api/", authorizeAjax(api))

	router.NotFound = http.HandlerFunc(requestCatchAll)
	log.Fatal(http.ListenAndServe(":9000", router))
}

func serveStaticFilesOr404Handler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	serveStaticFilesOr404(w, r)
}

func api(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	allowedFunctions := map[string]func(*http.Request) string{
		"notify-admin": apiNotifyAdmin,
		"tasks":        apiTasks,
		"new-task":     apiNewTask,
	}
	if r.PostFormValue("action") == "" {
		serve404(w)
		return
	}
	function, ok := allowedFunctions[r.PostFormValue("action")]
	if !ok {
		json.NewEncoder(w).Encode(ErrorJSON{
			Errors: []string{"The requested action could not be found in our api."},
		})
		return
	}
	fmt.Fprint(w, function(r))
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
