package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	// Assuming you have a function to register user routes
	// registerUserRoutes(r)
	return r
}

func TestUserRegistration(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		input      map[string]string
		statusCode int
	}{
		{"Valid User", map[string]string{"username": "testuser", "password": "password"}, http.StatusCreated},
		{"Missing Username", map[string]string{"password": "password"}, http.StatusBadRequest},
		{"Missing Password", map[string]string{"username": "testuser"}, http.StatusBadRequest},
		{"Empty Username", map[string]string{"username": "", "password": "password"}, http.StatusBadRequest},
		{"Empty Password", map[string]string{"username": "testuser", "password": ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestUserLogin(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		input      map[string]string
		statusCode int
	}{
		{"Valid Credentials", map[string]string{"username": "testuser", "password": "password"}, http.StatusOK},
		{"Invalid Credentials", map[string]string{"username": "wronguser", "password": "wrongpass"}, http.StatusUnauthorized},
		{"Missing Username", map[string]string{"password": "password"}, http.StatusBadRequest},
		{"Missing Password", map[string]string{"username": "testuser"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/users/testuser", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUserProfile(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		input      map[string]string
		statusCode int
	}{
		{"Valid Update", map[string]string{"bio": "Hello, I'm a test user!"}, http.StatusOK},
		{"Invalid Update", map[string]string{"bio": ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/api/users/testuser", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

// Additional tests for article CRUD operations and interactions can be added here.