package main

import (
	"fmt"
	"net/http"
	log "server/utils/login"
	todo "server/utils/todo"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Please login to gain access to the server.\n")
	fmt.Fprintf(w, "You can use http://host:port/Login?username={your_username}&password={your_password} to get the authentication tokens that can be passed in via Authorization header or as part of the POST body.\n")
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
