package main

import (
	"net/http"
	log "server/utils/login"
	todo "server/utils/todo"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.AuthMsg(w)
}

func initHttpRequests() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/Login", log.SignIn)

	http.HandleFunc("/List", todo.List)
	http.HandleFunc("/Add", todo.Add)
	http.HandleFunc("/Mark-complete", todo.Mark)
	http.HandleFunc("/Delete", todo.Delete)

	http.ListenAndServe(":8080", nil)
}

func main() {
	initHttpRequests()
}
