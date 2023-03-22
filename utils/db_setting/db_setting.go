package db_setting

// Configure database settings
func DbSettings() map[string]string {
	// Create a map with string keys and string values
	db := make(map[string]string)

	// Add some key-value pairs to the map
	db["user"] = "user"
	db["pw"] = "123456"
	db["conn"] = "tcp(127.0.0.1:3306)"
	db["schema"] = "todo"
	return db
}
