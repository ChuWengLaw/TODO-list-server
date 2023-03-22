package login

func SignIn(x int, y int) int {
	return x + y
}

// func userHandler(w http.ResponseWriter, r *http.Request) {
//     user := getCurrentUser(r)
//     if user == nil {
//         http.Redirect(w, r, "/login", http.StatusFound)
//         return
//     }

//     if !user.isAdmin {
//         http.Error(w, "Access denied", http.StatusForbidden)
//         return
//     }

//     fmt.Fprintf(w, "Welcome, %s!", user.name)
// }
