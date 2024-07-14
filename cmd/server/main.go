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
	// Запускаем наши миграции
	database.RunMigrations()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
	defer db.Close()

	http.HandleFunc("/api/auth/register", auth.RegisterHandler(db))
	http.HandleFunc("/api/auth/login", auth.LoginHandler(db))

	http.Handle("/api/groups", authMiddleware(groups.GetGroupsHandler(db)))
	http.Handle("/api/groups/create", authMiddleware(groups.CreateGroupHandler(db)))
	http.Handle("/api/getUser", authMiddleware(user.GetUserHandler(db)))

	log.Println("Connected to the database successfully")
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims, err := auth.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("username", claims.Username)
		next(w, r)
	}
}
