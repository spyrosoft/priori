package main

import (
	"errors"
	"net/http"
	"strconv"
)

func apiNewTask(w http.ResponseWriter, r *http.Request) interface{} {
	ok, response := validateNewTaskPostFormValues(r)
	if !ok {
		return response
	}

	newTask := r.PostFormValue("task")
	userId, _ := cookieIntValue("user-id", r)
	difficulty, _ := strconv.Atoi(r.PostFormValue("difficulty"))
	shortTerm, _ := strconv.Atoi(r.PostFormValue("short-term"))
	longTerm, _ := strconv.Atoi(r.PostFormValue("long-term"))

	var taskId int
	err := db.QueryRow("INSERT INTO user_tasks (user_id, task, difficulty, short_term, long_term) VALUES ($1, $2, $3, $4, $5) RETURNING id", userId, newTask, difficulty, shortTerm, longTerm).Scan(&taskId)
	if err != nil {
		return notifyAdminResponse("An error occurred while adding your new task to the database.", err)
	}
	return taskId
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
	userId, ok := currentUserId(r)
	if !ok {
		return notifyAdminResponse("An error occurred while looking up your tasks.", errors.New("Could not retrieve the current user id."))
	}

	isOrEqual := "is NULL"
	if r.PostFormValue("parent-id") != "" {
		parentId, err := strconv.Atoi(r.PostFormValue("parent-id"))
		if err != nil {
			return notifyAdminResponse("An error occurred accessing your tasks.", err)
		}
		safeParentId := strconv.Itoa(parentId)
		isOrEqual = "= " + safeParentId
	}
	var results interface{}
	err := db.QueryRow("SELECT TO_JSON(ARRAY_AGG(tasks)) FROM (SELECT id, task, difficulty, short_term, long_term FROM user_tasks WHERE user_id = $1 AND parent_id "+isOrEqual+") tasks", userId).Scan(&results)
	if err != nil {
		return notifyAdminResponse("An error occurred while looking up your tasks.", err)
	}

	switch results.(type) {
	case string:
		return results
	default:
		return apiResponse{Success: true}
	}
}

func apiDeleteTask(w http.ResponseWriter, r *http.Request) interface{} {
	unsafeId := r.PostFormValue("id")
	if unsafeId == "" {
		return apiResponse{Errors: []string{"Please provide the id for the task you would like to delete."}}
	}

	id, err := strconv.Atoi(unsafeId)
	userId, ok := currentUserId(r)
	if !ok {
		return notifyAdminResponse("An error occurred while looking up your user.", errors.New("Could not retrieve the current user id."))
	}

	//TODO: Cascade delete children
	_, err = db.Query("DELETE FROM user_tasks WHERE user_id = $1 AND id = $2", userId, id)
	if err != nil {
		return notifyAdminResponse("An error occurred while deleting your task.", err)
	}

	return apiResponse{Success: true}
}
