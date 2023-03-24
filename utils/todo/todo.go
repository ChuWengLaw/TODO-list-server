package todo

import (
	"database/sql"
	"fmt"
	"net/http"
	d "server/utils/db_setting"
	log "server/utils/login"

	_ "github.com/go-sql-driver/mysql"
)

// Helper function for DB connection
func DbConn() string {
	db_settings := fmt.Sprintf("%s:%s@%s/%s", d.DbSettings()["user"], d.DbSettings()["pw"], d.DbSettings()["conn"], d.DbSettings()["schema"])
	return db_settings
}

// List all user's tasks
func List(w http.ResponseWriter, r *http.Request) {
	// Check for curl header
	if log.CheckCurlHeader(w, r) {
		// Connect to db
		db, err := sql.Open("mysql", DbConn())

		if err != nil {
			fmt.Fprintf(w, "An error occured when connecting to database.")
			return
		}

		defer db.Close()

		// Query data
		slt_stmt, err := db.Prepare("SELECT * FROM todo_list WHERE user_id = ?")
		if err != nil {
			fmt.Fprintf(w, "An error occured when querying data.")
			return
		}

		defer slt_stmt.Close()

		rows, err := slt_stmt.Query(log.User_id)
		check_empty, err2 := slt_stmt.Query(log.User_id)
		if err != nil {
			fmt.Fprintf(w, "An error occured when executing statement.")
			return
		}
		if err2 != nil {
			fmt.Fprintf(w, "An error occured when executing statement.")
			return
		}

		if check_empty.Next() {
			// Extract data
			var id int
			var user_id int
			var todo string
			var is_completed int
			for rows.Next() {
				err := rows.Scan(&id, &user_id, &todo, &is_completed)
				if err != nil {
					fmt.Fprintf(w, "An error occured when scanning row.")
				}
				fmt.Fprintf(w, "User Id: %d, TODO: %s, Is Completed: %d\n", user_id, todo, is_completed)
			}
		} else {
			fmt.Fprintf(w, "Your TODO-list is empty, try adding a task by calling http://localhost:8080/Add?todo={your_task}")
		}
	} else {
		log.AuthMsg(w)
	}
}

// Add task to list
func Add(w http.ResponseWriter, r *http.Request) {
	// Check for curl header
	if log.CheckCurlHeader(w, r) {
		// Extract value of todo key from request param
		todo := r.URL.Query().Get("todo")

		// Check for empty value otherwise insert data to db
		if todo == "" {
			fmt.Fprintf(w, "Invalid task, try adding a task by invoking ?todo={your task} at the end of the api path.\n")
		} else {
			// Connect to db
			db, err := sql.Open("mysql", DbConn())

			if err != nil {
				fmt.Fprintf(w, "An error occured when connecting to database.")
				return
			}

			defer db.Close()

			// Prepare statement
			stmt, err := db.Prepare("INSERT INTO todo_list(user_id, todo, is_completed) VALUES(?, ?, ?)")
			if err != nil {
				fmt.Fprintf(w, "An error occured when preparing statement.")
				return
			}

			// Execute statement
			_, err = stmt.Exec(log.User_id, todo, 0)
			if err != nil {
				fmt.Fprintf(w, "An error occured when executing statement.")
				return
			}

			defer stmt.Close()

			fmt.Fprintf(w, "%s has been added to your TODO-list.\n", todo)
		}
	} else {
		log.AuthMsg(w)
	}
}

// Mark a task as completed in list
func Mark(w http.ResponseWriter, r *http.Request) {
	// Check for curl header
	if log.CheckCurlHeader(w, r) {
		// Extract value of todo key from request param
		todo := r.URL.Query().Get("todo")

		// Check for empty value otherwise edit data in db
		if todo == "" {
			fmt.Fprintf(w, "Invalid task, try adding a task by invoking ?todo={your task} at the end of the api path.\n")
		} else {
			// Connect to db
			db, err := sql.Open("mysql", DbConn())

			if err != nil {
				fmt.Fprintf(w, "An error occured when connecting to database.")
				return
			}

			defer db.Close()

			// Query data to check if the data exist then only we proceed to edit it
			slt_stmt, err := db.Prepare("SELECT COUNT(id) AS length FROM todo_list WHERE user_id = ? AND todo = ?")
			if err != nil {
				fmt.Fprintf(w, "An error occured when preparing statement.")
				return
			}

			defer slt_stmt.Close()

			rows, err := slt_stmt.Query(log.User_id, todo)
			if err != nil {
				fmt.Fprintf(w, "An error occured when executing statement.")
				return
			}

			defer rows.Close()

			var length int
			for rows.Next() {
				err := rows.Scan(&length)
				if err != nil {
					fmt.Fprintf(w, "An error occured when scanning row.")
				}
				if length > 0 {
					// Prepare statement
					stmt, err := db.Prepare("UPDATE todo_list SET is_completed = 1 WHERE user_id = ? AND todo = ?")
					if err != nil {
						fmt.Fprintf(w, "An error occured when preparing statement.")
					}

					// Execute statement
					_, err = stmt.Exec(log.User_id, todo)
					if err != nil {
						fmt.Fprintf(w, "An error occured when executing statement.")
					}

					defer stmt.Close()

					fmt.Fprintf(w, "%s in your TODO-list has been updated to 1.\n", todo)
				} else {
					fmt.Fprintf(w, "%s is not in your TODO-list.\n", todo)
				}
			}
		}
	} else {
		log.AuthMsg(w)
	}
}

// Delete task from list
func Delete(w http.ResponseWriter, r *http.Request) {
	// Check for curl header
	if log.CheckCurlHeader(w, r) {
		// Extract value of todo key from request param
		todo := r.URL.Query().Get("todo")

		// Check for empty value otherwise delete data to db
		if todo == "" {
			fmt.Fprintf(w, "Invalid task, try adding a task by invoking ?todo={your task} at the end of the api path.\n")
		} else {
			// Connect to db
			db, err := sql.Open("mysql", DbConn())

			if err != nil {
				fmt.Fprintf(w, "An error occured when connecting to database.")
				return
			}

			defer db.Close()

			// Query data to check if the data exist then only we proceed to delete it
			slt_stmt, err := db.Prepare("SELECT COUNT(id) AS length FROM todo_list WHERE user_id = ? AND todo = ?")
			if err != nil {
				fmt.Fprintf(w, "An error occured when preparing statement.")
				return
			}

			defer slt_stmt.Close()

			rows, err := slt_stmt.Query(log.User_id, todo)
			if err != nil {
				fmt.Fprintf(w, "An error occured when executing statement.")
				return
			}

			defer rows.Close()

			var length int
			for rows.Next() {
				err := rows.Scan(&length)
				if err != nil {
					fmt.Fprintf(w, "An error occured when scanning row.")
				}
				if length > 0 {
					// Prepare statement
					del_stmt, err := db.Prepare("DELETE FROM todo_list WHERE user_id = ? AND todo = ?")
					if err != nil {
						fmt.Fprintf(w, "An error occured when preparing statement.")
					}

					// Execute statement
					_, err = del_stmt.Exec(log.User_id, todo)
					if err != nil {
						fmt.Fprintf(w, "An error occured when executing statement.")
					}

					defer del_stmt.Close()

					fmt.Fprintf(w, "%s has been removed from your TODO-list.\n", todo)
				} else {
					fmt.Fprintf(w, "%s is not in your TODO-list.\n", todo)
				}
			}
		}
	} else {
		log.AuthMsg(w)
	}
}
