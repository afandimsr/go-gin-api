package postgres

import (
	"database/sql"
	"errors"

	"github.com/afandimsr/go-gin-api/internal/domain/user"
	"github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindAll(limit, offset int) ([]user.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) FindByID(id string) (user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	return u, nil
}

func (r *userRepo) Save(u user.User) error {
	id := uuid.New().String()

	_, err := r.db.Exec(
		"INSERT INTO users(id, name, email, password) VALUES($1, $2, $3, $4)",
		id, u.Name, u.Email, u.Password,
	)

	for _, role := range u.Roles {
		var roleID string
		err := r.db.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
		if err != nil {
			return err
		}
		_, err = r.db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)", id, roleID)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *userRepo) Update(u user.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		u.Name, u.Email, u.Password, u.ID,
	)
	return err
}

func (r *userRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepo) FindByEmail(email string) (user.User, error) {
	var u user.User

	query := `
		SELECT id, name, email, password, is_active
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(query, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}

	rows, err := r.db.Query(`
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, u.ID)

	if err != nil {
		return u, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return u, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return u, err
	}

	u.Roles = roles
	return u, nil
}
