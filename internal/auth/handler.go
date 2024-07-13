package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var creds UserCredentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		hashedPassword, err := HashPassword(creds.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, hashedPassword)
		if err != nil {
			http.Error(w, "Error saving user to database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	}
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		email := query["email"][0]
		password := query["password"][0]

		// var creds UserCredentials
		// err := json.NewDecoder(r.Body).Decode(&creds)
		// if err != nil {
		// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
		// 	return
		// }

		var storedPassword string
		err := db.QueryRow(`SELECT "Password" FROM "Users" WHERE "Email"=$1`, email).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				log.Printf("%w", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// if !CheckPasswordHash(creds.Password, storedPassword) {
		// 	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		// 	return
		// }

		if storedPassword != password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
	}
}
