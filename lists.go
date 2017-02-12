package main

import (
	"encoding/json"
	"net/http"
)

type NewListJSON struct {
	Success bool   `json:"success"`
	ID      int    `json:"name"`
	Name    string `json:"name"`
}

func newList(r *http.Request) (results string) {
	newListName := r.PostFormValue("new-list-name")
	if newListName == "" {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"The New List field is required."},
			Fields: []string{"new-list-name"},
		})
		results = string(resultsBytes)
		return
	}
	userID, _ := cookieIntValue("user-id", r)
	var listID int
	err := db.QueryRow("INSERT INTO user_lists (user_id, name) VALUES ($1, $2)", userID, newListName).Scan(&listID)
	if err != nil {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from adding the new list to the database."},
			Debug:  []string{err.Error()},
		})
		results = string(resultsBytes)
		return
	}
	resultsBytes, _ := json.Marshal(NewListJSON{
		Success: true,
		ID:      listID,
		Name:    newListName,
	})
	results = string(resultsBytes)
	return
}

func allLists(r *http.Request) (results string) {
	userID, ok := currentUserID(r)
	if !ok {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from looking up your lists."},
			Debug:  []string{"Could not retrieve the current user id."},
		})
		results = string(resultsBytes)
		return
	}
	err := db.QueryRow("SELECT TO_JSON(ARRAY_AGG(lists)) FROM (SELECT id, name FROM user_lists WHERE user_id = $1) lists", userID).Scan(results)
	if err != nil {
		resultsBytes, _ := json.Marshal(ErrorJSON{
			Errors: []string{"An error prevented us from looking up your lists."},
			Debug:  []string{err.Error()},
		})
		results = string(resultsBytes)
	}
	return
}
