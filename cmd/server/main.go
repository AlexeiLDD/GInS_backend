package main

import (
	"log"
	"net/http"

	"awesomeProject3/internal/auth"
	"awesomeProject3/internal/groups"
	"awesomeProject3/internal/user"
	"awesomeProject3/pkg/database"

	"github.com/rs/cors"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/getUser", auth.RegisterHandler(db))
	mux.HandleFunc("/api/auth/login", auth.LoginHandler(db))
	mux.HandleFunc("/api/groups", groups.GetGroupsHandler(db))
	mux.HandleFunc("/api/groups/create", groups.CreateGroupHandler(db))
	mux.HandleFunc("/api/getUser/{id}", user.GetUserHandler(db))

	c := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
			AllowedHeaders:   []string{"X-Requested-With", "Content-Type"},
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
			AllowCredentials: true,
			// Enable Debugging for testing, consider disabling in production
			Debug: true,
		},
	)

	handler := c.Handler(mux)

	log.Println("Starting server on :5062")
	log.Fatal(http.ListenAndServe(":5062", handler))
}
