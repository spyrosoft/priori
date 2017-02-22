package main

import (
	"errors"
	"net/http"
	"strconv"
)

type Task struct {
	Id         int    `json:"id"`
	Task       string `json:"task"`
	ShortTerm  int    `json:"short-term"`
	LongTerm   int    `json:"long-term"`
	Urgency    int    `json:"urgency"`
	Difficulty int    `json:"difficulty"`
}

func apiNewTask(w http.ResponseWriter, r *http.Request) interface{} {
	ok, response := validateNewTaskPostFormValues(r)
	if !ok {
		return response
	}

	newTask := r.PostFormValue("task")
	userId, _ := cookieIntValue("user-id", r)
	shortTerm, _ := strconv.Atoi(r.PostFormValue("short-term"))
	longTerm, _ := strconv.Atoi(r.PostFormValue("long-term"))
	urgency, _ := strconv.Atoi(r.PostFormValue("urgency"))
	difficulty, _ := strconv.Atoi(r.PostFormValue("difficulty"))

	var taskId int
	err := db.QueryRow("INSERT INTO user_tasks (user_id, task, short_term, long_term, urgency, difficulty) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", userId, newTask, shortTerm, longTerm, urgency, difficulty).Scan(&taskId)
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
	if _, err = strconv.Atoi(r.PostFormValue("short-term")); err != nil {
		response.Errors = append(response.Errors, "The Short Term field must be an integer.")
		response.Fields = append(response.Fields, "short-term")
	}
	if _, err = strconv.Atoi(r.PostFormValue("long-term")); err != nil {
		response.Errors = append(response.Errors, "The Long Term field must be an integer.")
		response.Fields = append(response.Fields, "long-term")
	}
	if _, err = strconv.Atoi(r.PostFormValue("urgency")); err != nil {
		response.Errors = append(response.Errors, "The Urgency field must be an integer.")
		response.Fields = append(response.Fields, "urgency")
	}
	if _, err = strconv.Atoi(r.PostFormValue("difficulty")); err != nil {
		response.Errors = append(response.Errors, "The Difficulty field must be an integer.")
		response.Fields = append(response.Fields, "difficulty")
	}
	return
}

func apiTasks(w http.ResponseWriter, r *http.Request) interface{} {
	userId, ok := currentUserId(r)
	if !ok {
		return notifyAdminResponse("An error occurred while looking up your tasks.", errors.New("Could not retrieve the current user id."))
	}

	isParentId := "is NULL"
	if r.PostFormValue("parent-id") != "" {
		parentId, err := strconv.Atoi(r.PostFormValue("parent-id"))
		if err != nil {
			return notifyAdminResponse("An error occurred accessing your tasks.", err)
		}
		safeParentId := strconv.Itoa(parentId)
		isParentId = "= " + safeParentId
	}
	rows, err := db.Query("SELECT id, task, short_term, long_term, urgency, difficulty FROM user_tasks WHERE user_id = $1 AND parent_id "+isParentId, userId)
	defer rows.Close()
	if err != nil {
		return notifyAdminResponse("An error occurred while looking up your tasks.", err)
	}

	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Task, &task.ShortTerm, &task.LongTerm, &task.Urgency, &task.Difficulty)
		if err != nil {
			return notifyAdminResponse("An error occurred while looking up your tasks.", err)
		}
		tasks = append(tasks, task)
	}
	return tasks
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
