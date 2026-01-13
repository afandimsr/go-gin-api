package mysql

import (
	"database/sql"
	"errors"

	"github.com/afandimsr/go-gin-api/internal/domain/user"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindAll(limit, offset int) ([]user.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users LIMIT ? OFFSET ?", limit, offset)
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

func (r *userRepo) FindByID(id int64) (user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	return u, nil
}

func (r *userRepo) Save(u user.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users(name, email, password) VALUES(?, ?, ?)",
		u.Name, u.Email, u.Password,
	)
	return err
}

func (r *userRepo) Update(u user.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		u.Name, u.Email, u.Password, u.ID,
	)
	return err
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *userRepo) FindByEmail(email string) (user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	return u, nil
}
