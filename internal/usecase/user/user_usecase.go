package user

import (
	"errors"
	"strings"

	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
	"github.com/afandimsr/go-gin-api/internal/domain/valueobject"
	pw "github.com/afandimsr/go-gin-api/internal/domain/valueobject"
	"github.com/afandimsr/go-gin-api/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	repo        user.UserRepository
	authService user.AuthService
}

func New(repo user.UserRepository, authService user.AuthService) *Usecase {
	return &Usecase{
		repo:        repo,
		authService: authService,
	}
}

func (u *Usecase) GetAll(page, limit int) ([]user.User, error) {
	offset := (page - 1) * limit
	return u.repo.FindAll(limit, offset)
}

func (u *Usecase) GetByID(id string) (user.User, error) {
	availableUser, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return user.User{}, apperror.NotFound(
				"User tidak ditemukan",
				err,
			).WithCode(apperror.UserNotFound)
		}

		return user.User{}, apperror.Internal(err)
	}

	return availableUser, nil
}

func (u *Usecase) Create(newUser user.User) error {
	if newUser.Email == "" {
		return apperror.BadRequest("email is required", nil)
	}

	defaultPassord := strings.Split(newUser.Email, "@")[0] + "123"
	if newUser.Password == "" {
		newUser.Password = defaultPassord
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Internal(err)
	}
	newUser.Password = string(hashedPassword)

	if err := u.repo.Save(newUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Update(id string, updatedUser user.User) error {
	if updatedUser.Email == "" {
		return apperror.BadRequest("email is required", nil)
	}

	// Check if user exists
	existingUser, err := u.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return apperror.NotFound(
				"User tidak ditemukan",
				err,
			).WithCode(apperror.UserNotFound)
		}

		return apperror.Internal(err)
	}

	existingUser.Name = updatedUser.Name
	existingUser.Email = updatedUser.Email
	existingUser.Roles = updatedUser.Roles

	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperror.Internal(err)
		}
		existingUser.Password = string(hashedPassword)
	}

	if err := u.repo.Update(existingUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Delete(id string) error {
	// Check if user exists
	if _, err := u.repo.FindByID(id); err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return apperror.NotFound(
				"User tidak ditemukan",
				err,
			).WithCode(apperror.UserNotFound)
		}

		return apperror.Internal(err)
	}

	if err := u.repo.Delete(id); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

func (u *Usecase) Login(email, password string) (string, error) {
	// 1. Find user by email
	existingUser, err := u.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", apperror.Unauthorized(
				"Username/Password tidak valid!",
				err,
			).WithCode(apperror.InvalidCredentials)
		}

		return "", apperror.Internal(err)
	}

	// 2. Authenticate
	// If AuthService is available (and configured), try it first/instead.
	// For this implementation, if AuthService is provided, we use it to validate password.
	// If it returns true, we consider it valid.
	authenticated := false
	if u.authService != nil {
		isAuth, err := u.authService.Login(email, password)
		if err == nil && isAuth {
			authenticated = true
		}
	}

	// Fallback to local bcrypt if not authenticated via external service (or if service not used)
	// Note: The requirement implies "if login using client auth service".
	// We can assume priority: External > Local.
	if !authenticated {
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
			return "", apperror.Unauthorized("Username/Password tidak valid!", nil).WithCode(apperror.InvalidCredentials)
		}
	}

	// 3. Generate Token
	token, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
	if err != nil {
		return "", apperror.Internal(err)
	}

	return token, nil
}

// ChangePassword changes the password of a user
func (u *Usecase) ChangePassword(id string, newPassword string) error {
	// Add validation password
	pw, err := pw.Password(newPassword)

	if err != nil {
		switch err {
		case valueobject.ErrPasswordTooShort:
			return apperror.Validation(err).
				WithCode(apperror.PasswordTooShort)

		case valueobject.ErrPasswordNoUpper,
			valueobject.ErrPasswordNoLower,
			valueobject.ErrPasswordNoDigit,
			valueobject.ErrPasswordNoSpecial:
			return apperror.Validation(err).
				WithCode(apperror.PasswordWeak)
		}
	}

	// Check if user exists
	if _, err := u.repo.FindByID(id); err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return apperror.NotFound(
				"User tidak ditemukan",
				err,
			).WithCode(apperror.UserNotFound)
		}

		return apperror.Internal(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(pw),
		bcrypt.DefaultCost)

	if err != nil {
		return apperror.Internal(err)
	}

	if err := u.repo.ChangePassword(id, string(hashedPassword)); err != nil {
		return apperror.Internal(err)
	}

	return nil
}
