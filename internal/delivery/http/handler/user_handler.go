package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/go-gin-api/internal/delivery/http/response"
	"github.com/username/go-gin-api/internal/domain/apperror"
	"github.com/username/go-gin-api/internal/domain/user"
	uc "github.com/username/go-gin-api/internal/usecase/user"
)

type UserHandler struct {
	usecase *uc.Usecase
}

func New(usecase *uc.Usecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

// GetUsers godoc
// @Summary      Get all users
// @Tags         Users
// @Produce      json
// @Success      200 {object} response.SuccessUserResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.usecase.GetAll(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_FETCH_ERROR", "failed to get users", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "success", users)
}

// GetUser godoc
// @Summary      Get user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	u, err := h.usecase.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", u)
}

// CreateUser godoc
// @Summary      Create user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body user.User true "User payload"
// @Success      201 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	if err := h.usecase.Create(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "user created", nil)
}

// Login godoc
// @Summary      Login user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body user.LoginRequest true "Login payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "login success", gin.H{"token": token})
}

// UpdateUser godoc
// @Summary      Update user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        body body user.User true "User payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	var req user.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	if err := h.usecase.Update(id, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user updated", nil)
}

// DeleteUser godoc
// @Summary      Delete user
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user deleted", nil)
}
