package mysql

import (
	"database/sql"

	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
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
	rows, err := r.db.Query("SELECT id, name, email, is_active FROM users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, apperror.HandleDatabaseError(err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive)
		//fetch roles for user
		roleRows, err := r.db.Query(`
			SELECT r.name
			FROM roles r
			JOIN user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = ?
		`, u.ID)
		if err != nil {
			return nil, apperror.HandleDatabaseError(err)
		}
		defer roleRows.Close()

		var roles []string
		for roleRows.Next() {
			var role string
			if err := roleRows.Scan(&role); err != nil {
				return nil, err
			}
			roles = append(roles, role)
		}

		u.Roles = roles
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) FindByID(id string) (user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, user.ErrUserNotFound
		}
		return u, apperror.HandleDatabaseError(err)
	}
	return u, nil
}

func (r *userRepo) Save(u user.User) error {
	id := uuid.New().String()
	_, err := r.db.Exec(
		"INSERT INTO users(id, name, email, password) VALUES(?, ?, ?, ?)",
		id, u.Name, u.Email, u.Password,
	)

	for _, role := range u.Roles {
		var roleID string
		err := r.db.QueryRow("SELECT id FROM roles WHERE name = ?", role).Scan(&roleID)
		if err != nil {
			return err
		}
		_, err = r.db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)", id, roleID)
		if err != nil {
			return apperror.HandleDatabaseError(err)
		}
	}

	return apperror.HandleDatabaseError(err)
}

func (r *userRepo) Update(u user.User) error {

	_, err := r.db.Exec(
		"UPDATE users SET name = ?, email = ?, password = ? WHERE id = ?",
		u.Name, u.Email, u.Password, u.ID,
	)

	// Update roles
	_, err = r.db.Exec("DELETE FROM user_roles WHERE user_id = ?", u.ID)
	if err != nil {
		return err
	}

	for _, role := range u.Roles {
		var roleID string
		err := r.db.QueryRow("SELECT id FROM roles WHERE name = ?", role).Scan(&roleID)
		if err != nil {
			return err
		}
		_, err = r.db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)", u.ID, roleID)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *userRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return apperror.HandleDatabaseError(err)
}

func (r *userRepo) FindByEmail(email string) (user.User, error) {
	var u user.User
	err := r.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, user.ErrUserNotFound
		}
		return u, apperror.HandleDatabaseError(err)
	}

	rows, err := r.db.Query(`
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?
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
		return u, apperror.HandleDatabaseError(err)
	}

	u.Roles = roles
	return u, nil
}

func (r *userRepo) ChangePassword(id string, newPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = ? WHERE id = ?", newPassword, id)
	return apperror.HandleDatabaseError(err)
}
