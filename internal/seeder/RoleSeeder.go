package seeder

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

func SeedRoles(db *sql.DB) {

	roles := []string{"ADMIN", "USER"}

	for _, role := range roles {
		id := uuid.NewString()
		err := db.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert role
				_, err = db.Exec("INSERT INTO roles(id,name) VALUES($1,$2)", id, role)
				if err != nil {
					log.Fatalf("failed to insert role %s: %v", role, err)
				}
			}
		} else {
			log.Printf("Role %s already exists with id %s", role, id)
		}
	}
}
