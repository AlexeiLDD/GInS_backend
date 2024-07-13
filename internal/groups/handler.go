package groups

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var group Group
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO groups (name, description) VALUES ($1, $2)", group.Name, group.Description)
		if err != nil {
			http.Error(w, "Error creating group", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Group created successfully"))
	}
}

func GetGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, description FROM groups")
		if err != nil {
			http.Error(w, "Error retrieving groups", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var groups []Group
		for rows.Next() {
			var group Group
			if err := rows.Scan(&group.ID, &group.Name, &group.Description); err != nil {
				http.Error(w, "Error scanning group", http.StatusInternalServerError)
				return
			}
			groups = append(groups, group)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error with rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
