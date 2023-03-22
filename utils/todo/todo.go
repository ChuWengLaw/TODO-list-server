package todo

import (
	"fmt"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List.")
}

func Add(w http.ResponseWriter, r *http.Request) {
	todo := r.URL.Query().Get("todo")
	if todo == "" {
		fmt.Fprintf(w, "New task not found, try adding a task by invoking ?todo={your task} at the end of the api path.\n")
	} else {
		fmt.Fprintf(w, "%s has been added to your TODO-list.\n", todo)
	}
}

func Mark(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Mark.")
}

func Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete.")
}
