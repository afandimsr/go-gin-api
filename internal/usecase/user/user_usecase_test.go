package user_test

import (
	"errors"
	"testing"

	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
	uc "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of user.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(limit, offset int) ([]user.User, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (user.User, error) {
	args := m.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) Save(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) ChangePassword(id string, newPassword string) error {
	args := m.Called(id, newPassword)
	return args.Error(0)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	mockUser := user.User{ID: "ef6d1df7-f85c-426c-9c12-6d58a1fc2633", Name: "Test User", Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", "ef6d1df7-f85c-426c-9c12-6d58a1fc2633").Return(mockUser, nil)

		u, err := usecase.GetByID("ef6d1df7-f85c-426c-9c12-6d58a1fc2633")

		assert.NoError(t, err)
		assert.Equal(t, mockUser.ID, u.ID)
		assert.Equal(t, mockUser.Name, u.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("FindByID", "412fd1c1-0d29-45dd-9cb7-efcf64390e8b").Return(user.User{}, user.ErrUserNotFound)

		_, err := usecase.GetByID("412fd1c1-0d29-45dd-9cb7-efcf64390e8b")

		assert.Error(t, err)
		assert.Equal(t, "User tidak ditemukan", err.Error())
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		newUser := user.User{Name: "New User", Email: "new@example.com", Password: "password123"}

		mockRepo.On("Save", mock.AnythingOfType("user.User")).Return(nil).Once()
		err := usecase.Create(newUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingEmail", func(t *testing.T) {
		err := usecase.Create(user.User{Password: "123"})
		assert.Error(t, err)
	})
}

func TestChangePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)
	t.Run("Success", func(t *testing.T) {
		userID := "6906ab46-7eda-4df8-8ad4-f9b46e39cb32"
		newPassword := "Newpassword123@"

		// ✅ mock FindByID (WAJIB)
		mockRepo.
			On("FindByID", userID).
			Return(user.User{ID: userID}, nil).
			Once()

		// ✅ mock ChangePassword
		mockRepo.
			On("ChangePassword", userID, mock.MatchedBy(func(pw string) bool {
				return pw != newPassword && len(pw) > 20
			})).
			Return(nil)

		err := usecase.ChangePassword(userID, newPassword)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := "non-existent-id"
		newPassword := "Newpassword123@"

		mockRepo.
			On("FindByID", userID).
			Return(user.User{}, user.ErrUserNotFound)

		err := usecase.ChangePassword(userID, newPassword)

		assert.Error(t, err)
		assert.Equal(t, "User tidak ditemukan", err.Error())

		var appErr *apperror.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperror.UserNotFound, appErr.ErrorCode)
	})

	t.Run("WeakPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := uc.New(mockRepo, nil)

		// ✅ mock FindByID (WAJIB)
		mockRepo.
			On("FindByID", "some-id").
			Return(user.User{ID: "some-id"}, nil).
			Once()

		err := usecase.ChangePassword("some-id", "newpassword123")

		assert.Error(t, err)

		var appErr *apperror.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperror.PasswordWeak, appErr.ErrorCode)
	})

	t.Run("ShortPassword", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		usecase := uc.New(mockRepo, nil)

		// ✅ mock FindByID (WAJIB)
		mockRepo.
			On("FindByID", "some-id").
			Return(user.User{ID: "some-id"}, nil).
			Once()

		err := usecase.ChangePassword("some-id", "123")

		assert.Error(t, err)

		var appErr *apperror.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperror.PasswordTooShort, appErr.ErrorCode)
	})

}

func TestDelete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)
	t.Run("Success", func(t *testing.T) {
		userID := "6906ab46-7eda-4df8-8ad4-f9b46e39cb32"
		// ✅ mock FindByID (WAJIB)
		mockRepo.
			On("FindByID", userID).
			Return(user.User{ID: userID}, nil).
			Once()
		// ✅ mock Delete
		mockRepo.
			On("Delete", userID).
			Return(nil).
			Once()
		err := usecase.Delete(userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("UserNotFound", func(t *testing.T) {
		userID := "non-existent-id"
		mockRepo.
			On("FindByID", userID).
			Return(user.User{}, user.ErrUserNotFound)
		err := usecase.Delete(userID)

		assert.Error(t, err)
		assert.Equal(t, "User tidak ditemukan", err.Error())
		var appErr *apperror.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperror.UserNotFound, appErr.ErrorCode)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)
	t.Run("Success", func(t *testing.T) {
		userID := "6906ab46-7eda-4df8-8ad4-f9b46e39cb32"
		updatedUser := user.User{Name: "Updated User", Email: "updated@example.com", Roles: []string{"USER"}, Password: "newpassword123"}
		// ✅ mock FindByID (WAJIB)
		mockRepo.
			On("FindByID", userID).
			Return(user.User{ID: userID, Name: "Old Name", Email: "old@example.com", Roles: []string{"USER"}, Password: "oldpassword123"}, nil).
			Once()

		// ✅ mock Update
		mockRepo.
			On("Update", mock.AnythingOfType("user.User")).
			Return(nil).
			Once()

		err := usecase.Update(userID, updatedUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("UserNotFound", func(t *testing.T) {
		userID := "non-existent-id"
		updatedUser := user.User{Name: "Updated User", Email: "updated@example.com"}
		mockRepo.
			On("FindByID", userID).
			Return(user.User{}, user.ErrUserNotFound)
		err := usecase.Update(userID, updatedUser)
		assert.Error(t, err)
		assert.Equal(t, "User tidak ditemukan", err.Error())
		var appErr *apperror.AppError
		assert.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperror.UserNotFound, appErr.ErrorCode)
	})
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)
	t.Run("Success", func(t *testing.T) {
		mockUsers := []user.User{
			{ID: "1", Name: "User One", Email: "user1@example.com"},
			{ID: "2", Name: "User Two", Email: "user2@example.com"},
		}

		mockRepo.
			On("FindAll", 10, 0).
			Return(mockUsers, nil).
			Once()
		users, err := usecase.GetAll(1, 10)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyResult", func(t *testing.T) {
		mockRepo.
			On("FindAll", 10, 0).
			Return([]user.User{}, nil).
			Once()
		users, err := usecase.GetAll(1, 10)
		assert.NoError(t, err)
		assert.Len(t, users, 0)
		mockRepo.AssertExpectations(t)
	})

}
