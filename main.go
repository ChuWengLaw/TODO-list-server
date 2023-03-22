package main

import (
	"fmt"
	"net/http"
	log "server/utils/login"
	todo "server/utils/todo"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func initHttpRequests() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/List", todo.List)
	http.HandleFunc("/Add", todo.Add)
	http.HandleFunc("/Mark-complete", todo.Mark)
	http.HandleFunc("/Delete", todo.Delete)
	http.ListenAndServe(":8080", nil)
}

func main() {
	test := log.SignIn(1, 2)
	fmt.Println(test)

	initHttpRequests()
}
