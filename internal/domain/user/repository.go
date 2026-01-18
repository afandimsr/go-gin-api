package user

type UserRepository interface {
	FindAll(limit, offset int) ([]User, error)
	FindByID(id string) (User, error)
	FindByEmail(email string) (User, error)
	Save(user User) error
	Update(user User) error
	Delete(id string) error
	ChangePassword(id string, newPassword string) error // New method for changing password
}

type AuthService interface {
	Login(email, password string) (bool, error)
}
