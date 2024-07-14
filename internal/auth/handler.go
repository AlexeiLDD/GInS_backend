package auth

import (
	"awesomeProject3/internal/server"
	"database/sql"
	"encoding/json"
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

		_, err = db.Exec(
			`
			INSERT INTO "Users" ("Id", "Password", "Email", "Name", "Telegram", "Phone", "AvatarId") 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			`,
			creds.ID, hashedPassword, creds.Email, creds.Username, creds.Telegram, creds.Phone, creds.AvatarId,
		)
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

		var storedPassword, storedId, storedName string
		err := db.QueryRow(`SELECT "Password", "Id", "Name" FROM "Users" WHERE "Email"=$1`, email).
			Scan(&storedPassword, &storedId, &storedName)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		if !CheckPasswordHash(password, storedPassword) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		server.SendOkResponse(w, LoginResponce{
			ID:       storedId,
			Email:    email,
			Username: storedName,
		})
	}
}
