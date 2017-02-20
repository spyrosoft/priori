package main

import (
	"net/http"
	"strconv"
)

func apiNewTask(w http.ResponseWriter, r *http.Request) interface{} {
	ok, response := validateNewTaskPostFormValues(r)
	if !ok {
		return response
	}

	newTask := r.PostFormValue("task")
	userID, _ := cookieIntValue("user-id", r)
	difficulty, _ := strconv.Atoi(r.PostFormValue("difficulty"))
	shortTerm, _ := strconv.Atoi(r.PostFormValue("short-term"))
	longTerm, _ := strconv.Atoi(r.PostFormValue("long-term"))

	var taskID int
	err := db.QueryRow("INSERT INTO user_tasks (user_id, task, difficulty, short_term, long_term) VALUES ($1, $2, $3, $4, $5) RETURNING id", userID, newTask, difficulty, shortTerm, longTerm).Scan(&taskID)
	if err != nil {
		return notifyAdminResponse("An error occurred while adding your new task to the database.", err)
	}
	return userID
}

func validateNewTaskPostFormValues(r *http.Request) (
	ok bool,
	response apiResponse,
) {
	ok = true
	var err error
	task := r.PostFormValue("task")
	if task == "" {
		ok = false
		response.Errors = append(response.Errors, "The Task field is required.")
		response.Fields = append(response.Fields, "task")
	}
	if _, err = strconv.Atoi(r.PostFormValue("difficulty")); err != nil {
		response.Errors = append(response.Errors, "The Difficulty field must be an integer.")
		response.Fields = append(response.Fields, "difficulty")
	}
	if _, err = strconv.Atoi(r.PostFormValue("short-term")); err != nil {
		response.Errors = append(response.Errors, "The Short Term field must be an integer.")
		response.Fields = append(response.Fields, "short-term")
	}
	if _, err = strconv.Atoi(r.PostFormValue("long-term")); err != nil {
		response.Errors = append(response.Errors, "The Long Term field must be an integer.")
		response.Fields = append(response.Fields, "long-term")
	}
	return
}

func apiTasks(w http.ResponseWriter, r *http.Request) interface{} {
	userID, ok := currentUserID(r)
	if !ok {
		return apiResponse{
			Errors: []string{"An error prevented us from looking up your tasks."},
			Debug:  []string{"Could not retrieve the current user id."},
		}
	}

	var results string
	err := db.QueryRow("SELECT TO_JSON(ARRAY_AGG(tasks)) FROM (SELECT id, task, difficulty, short_term, long_term FROM user_tasks WHERE user_id = $1) tasks", userID).Scan(&results)
	if err != nil {
		return apiResponse{
			Errors: []string{"An error prevented us from looking up your tasks."},
			Debug:  []string{err.Error()},
		}
	}
	return results
}
