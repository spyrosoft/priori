package main

import (
	"encoding/json"
	"net/http"
)

type NewTaskJSON struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

func apiNewTask(r *http.Request) (results string) {
	newTask := r.PostFormValue("task")
	if newTask == "" {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"The Task field is required."},
			Fields: []string{"task"},
		})
		results = string(resultsBytes)
		return
	}
	userID, _ := cookieIntValue("user-id", r)
	var taskID int
	err := db.QueryRow("INSERT INTO user_tasks (user_id, task) VALUES ($1, $2) RETURNING id", userID, newTask).Scan(&taskID)
	if err != nil {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from adding the new task to the database."},
			Debug:  []string{err.Error()},
		})
		results = string(resultsBytes)
		return
	}
	resultsBytes, _ := json.Marshal(NewTaskJSON{
		ID:   taskID,
		Task: newTask,
	})
	results = string(resultsBytes)
	return
}

func apiTasks(r *http.Request) (results string) {
	userID, ok := currentUserID(r)
	if !ok {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from looking up your tasks."},
			Debug:  []string{"Could not retrieve the current user id."},
		})
		results = string(resultsBytes)
		return
	}
	err := db.QueryRow("SELECT TO_JSON(ARRAY_AGG(tasks)) FROM (SELECT id, task FROM user_tasks WHERE user_id = $1) tasks", userID).Scan(&results)
	if err != nil {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from looking up your tasks."},
			Debug:  []string{err.Error()},
		})
		results = string(resultsBytes)
	}
	return
}
