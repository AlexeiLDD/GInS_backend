package main

import (
	"log"
	"net/http"

	"awesomeProject3/internal/auth"
	"awesomeProject3/internal/groups"
	"awesomeProject3/internal/user"
	"awesomeProject3/pkg/database"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
	defer db.Close()

	http.HandleFunc("/api/auth/register", auth.RegisterHandler(db))
	http.HandleFunc("/api/auth/login", auth.LoginHandler(db))
	http.HandleFunc("/api/groups", groups.GetGroupsHandler(db))
	http.HandleFunc("/api/groups/create", groups.CreateGroupHandler(db))
	http.HandleFunc("/api/getUser", user.GetUserHandler(db))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
