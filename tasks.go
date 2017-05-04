package main

import (
	"database/sql"
	"errors"
	"net/http"
)

type Tasks struct {
	Parent   string `json:"parent,omitempty"`
	ParentId int    `json:"parent-id,omitempty"`
	Tasks    string `json:"tasks"`
}

func apiGetTasks(w http.ResponseWriter, r *http.Request) interface{} {
	generalUserError := "An error occurred while looking up your tasks."
	userId, ok := currentUserId(r)
	if !ok {
		return notifyAdminResponse(generalUserError, errors.New("Could not retrieve the current user id."))
	}

	var tasks string
	var s sql.NullString

	err := db.QueryRow("SELECT tasks FROM users WHERE id = $1", userId).Scan(&s)
	if err != nil {
		return notifyAdminResponse(generalUserError, err)
	}
	if s.Valid {
		tasks = s.String
	}

	if tasks == "" {
		return "[]"
	}

	return tasks
}

func apiUpdateTasks(w http.ResponseWriter, r *http.Request) interface{} {
	tasks := r.PostFormValue("tasks")
	userId, _ := cookieIntValue("user-id", r)

	//TODO: Validate this before releasing it into the wild
	_, err := db.Query("UPDATE users SET tasks = $1 WHERE id = $2", tasks, userId)
	if err != nil {
		return notifyAdminResponse("An error occurred while updating your tasks in the database.", err)
	}
	return apiResponse{Success: true}
}
