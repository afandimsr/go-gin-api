package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user/request"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/helper"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/response"
	"github.com/afandimsr/go-gin-api/internal/domain/apperror"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
	"github.com/afandimsr/go-gin-api/internal/pkg/oidc"
	uc "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase      *uc.Usecase
	oidcProvider *oidc.OIDCProvider
}

func New(usecase *uc.Usecase, oidcProvider *oidc.OIDCProvider) *UserHandler {
	return &UserHandler{
		usecase:      usecase,
		oidcProvider: oidcProvider,
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
		c.Error(err)
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
	id, err := helper.ValidateUUIDParamNotFound(c, "id")
	if err != nil {
		c.Error(err)
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
		c.Error(apperror.Validation(err).WithCode(apperror.ValidationError))
		return
	}

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "login success", user.LoginResponse{Token: token})
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
	id, err := helper.ValidateUUIDParamNotFound(c, "id")
	if err != nil {
		c.Error(err)
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
	id, err := helper.ValidateUUIDParamNotFound(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user deleted", nil)
}

// ChangePassword godoc
// @Summary      Change user password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        body body request.ChangePasswordRequest true "Change password payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id}/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, err := helper.ValidateUUIDParamNotFound(c, "id")
	if err != nil {
		c.Error(err)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(err).WithCode(apperror.ValidationError))
		return
	}

	if err := h.usecase.ChangePassword(id, req.NewPassword); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "password changed", nil)
}

func (h *UserHandler) OIDCLogin(c *gin.Context) {
	if h.oidcProvider == nil {
		c.Error(apperror.Internal(fmt.Errorf("OIDC provider not configured")))
		return
	}

	state := "random-state" // Should be dynamic and validated
	url := h.oidcProvider.OAuth2Config.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (h *UserHandler) OIDCCallback(c *gin.Context) {
	if h.oidcProvider == nil {
		c.Error(apperror.Internal(fmt.Errorf("OIDC provider not configured")))
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	log.Printf("[OIDC] Callback received. Code: %s, State: %s", code, state)

	if state != "random-state" {
		log.Printf("[OIDC] State mismatch. Expected: random-state, Got: %s", state)
		c.Error(apperror.BadRequest("invalid state", nil))
		return
	}

	oauth2Token, err := h.oidcProvider.OAuth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Printf("[OIDC] Token exchange failed: %v", err)
		c.Error(apperror.Internal(err))
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Printf("[OIDC] No id_token in response")
		c.Error(apperror.Internal(fmt.Errorf("no id_token in token response")))
		return
	}

	idToken, err := h.oidcProvider.Verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		log.Printf("[OIDC] ID Token verification failed: %v", err)
		c.Error(apperror.Unauthorized("failed to verify ID token", err))
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("[OIDC] Failed to extract claims: %v", err)
		c.Error(apperror.Internal(err))
		return
	}

	log.Printf("[OIDC] Claims extracted: email=%v, sub=%v", claims["email"], claims["sub"])

	token, err := h.usecase.LoginWithOIDC(claims, oauth2Token.AccessToken)
	if err != nil {
		log.Printf("[OIDC] LoginWithOIDC failed: %v", err)
		c.Error(err)
		return
	}

	log.Printf("[OIDC] Login successful, redirecting to frontend")

	// Redirect to frontend with token
	targetURL := fmt.Sprintf("http://localhost:5173/auth/callback?token=%s", token)
	c.Redirect(http.StatusFound, targetURL)
}

func (h *UserHandler) Logout(c *gin.Context) {
	if h.oidcProvider == nil {
		response.Success(c, http.StatusOK, "logout success (local)", nil)
		return
	}

	issuer := h.oidcProvider.IssuerURL
	// Redirect back to frontend login page after Keycloak logout
	postLogoutRedirect := "http://localhost:5173/login"
	logoutURL := fmt.Sprintf("%s/protocol/openid-connect/logout?client_id=%s&post_logout_redirect_uri=%s",
		issuer, h.oidcProvider.OAuth2Config.ClientID, postLogoutRedirect)

	c.Redirect(http.StatusFound, logoutURL)
}
