package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
)

type keycloakService struct {
	cfg    config.KeycloakConfig
	client *http.Client
}

func NewKeycloakService(cfg config.KeycloakConfig) user.KeycloakService {
	return &keycloakService{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *keycloakService) getAdminToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", "admin-cli")
	data.Set("username", s.cfg.AdminUser)
	data.Set("password", s.cfg.AdminPassword)
	data.Set("grant_type", "password")

	u := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", s.cfg.URL)
	resp, err := s.client.PostForm(u, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get admin token: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (s *keycloakService) CreateUser(email, name, password string, roles []string) (string, error) {
	token, err := s.getAdminToken()
	if err != nil {
		return "", err
	}

	userData := map[string]interface{}{
		"username":      email,
		"email":         email,
		"enabled":       true,
		"emailVerified": true,
		"firstName":     name,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	body, _ := json.Marshal(userData)
	u := fmt.Sprintf("%s/admin/realms/%s/users", s.cfg.URL, s.cfg.Realm)
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		// User already exists, let's try to find their ID
		return s.findUserIDByEmail(email, token)
	}

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create user: %s, body: %s", resp.Status, string(b))
	}

	// Keycloak returns user ID in Location header
	location := resp.Header.Get("Location")
	parts := strings.Split(location, "/")
	return parts[len(parts)-1], nil
}

func (s *keycloakService) findUserIDByEmail(email, token string) (string, error) {
	u := fmt.Sprintf("%s/admin/realms/%s/users?username=%s", s.cfg.URL, s.cfg.Realm, email)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var users []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", fmt.Errorf("user not found after conflict")
	}

	return users[0].ID, nil
}

func (s *keycloakService) VerifyToken(accessToken string) error {
	u := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", s.cfg.URL, s.cfg.Realm)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("session invalid (Keycloak returned %d)", resp.StatusCode)
	}

	return nil
}
