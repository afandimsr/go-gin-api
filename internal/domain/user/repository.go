package user

type UserRepository interface {
	FindAll(limit, offset int) ([]User, error)
	FindByID(id string) (User, error)
	FindByEmail(email string) (User, error)
	FindByKeycloakID(keycloakID string) (User, error)
	Save(user User) error
	Update(user User) error
	UpdateKeycloakID(id string, keycloakID string) error
	Delete(id string) error
	ChangePassword(id string, newPassword string) error // New method for changing password
}

type AuthService interface {
	Login(email, password string) (bool, error)
}

type KeycloakService interface {
	CreateUser(email, name, password string, roles []string) (string, error)
	VerifyToken(accessToken string) error
}
