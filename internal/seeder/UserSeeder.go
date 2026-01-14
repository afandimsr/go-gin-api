package seeder

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdminUser seeds a default admin user into the database.
func SeedAdminUserPostgres(db *sql.DB) {
	id := uuid.NewString()
	name := "Admin"
	email := "admin@example.com"
	password := "admin123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var userID string
	// Check if user exists
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert user
			err = db.QueryRow("INSERT INTO users(id,name, email, password, is_active) VALUES($1, $2, $3, $4, $5) RETURNING id", id, name, email, string(hashed), true).Scan(&userID)
			if err != nil {
				log.Fatalf("failed to insert user: %v", err)
			}
		} else {
			log.Fatalf("failed to query user: %v", err)
		}
	}

	// Find ADMIN role id
	var roleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = $1", "ADMIN").Scan(&roleID)
	if err != nil {
		log.Fatalf("failed to find ADMIN role: %v", err)
	}

	// Insert into user_roles
	_, err = db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES($1, $2) ON CONFLICT DO NOTHING", userID, roleID)
	if err != nil {
		log.Fatalf("failed to insert user_roles: %v", err)
	}

	log.Printf("Admin user seeded: id=%s, email=%s", id, email)
}

func SeedAdminUserMysql(db *sql.DB) {
	id := uuid.NewString()
	name := "Admin"
	email := "admin@example.com"
	password := "admin123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var userID string
	// Check if user exists
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert user
			_, err := db.Exec(
				"INSERT INTO users(id,name, email, password, is_active) VALUES(?, ?, ?, ?, ?)",
				id, name, email, string(hashed), true,
			)

			if err != nil {
				log.Fatalf("failed to insert user: %v", err)
			}

		} else {
			log.Fatalf("failed to query user: %v", err)
		}
	}

	// Find ADMIN role id
	var roleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = ?", "ADMIN").Scan(&roleID)
	if err != nil {
		log.Fatalf("failed to find ADMIN role: %v", err)
	}

	// Insert into user_roles
	_, err = db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)", id, roleID)
	if err != nil {
		log.Fatalf("failed to insert user_roles: %v", err)
	}

	log.Printf("Admin user seeded: id=%s, email=%s", id, email)
}
