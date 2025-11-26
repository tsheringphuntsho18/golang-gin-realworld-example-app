package main

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
	// Define your routes here
	return r
}

func TestUserAuthentication(t *testing.T) {
	r := setupRouter()

	// Test user registration
	t.Run("User Registration", func(t *testing.T) {
		user := map[string]string{"username": "testuser", "password": "password"}
		jsonUser, _ := json.Marshal(user)

		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test user login
	t.Run("User Login", func(t *testing.T) {
		user := map[string]string{"username": "testuser", "password": "password"}
		jsonUser, _ := json.Marshal(user)

		req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(jsonUser))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test user profile retrieval
	t.Run("Get User Profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/users/testuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestArticleCRUD(t *testing.T) {
	r := setupRouter()

	// Test article creation
	t.Run("Create Article", func(t *testing.T) {
		article := map[string]string{"title": "Test Article", "body": "This is a test article."}
		jsonArticle, _ := json.Marshal(article)

		req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(jsonArticle))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test article retrieval
	t.Run("Get Article", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/articles/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test article update
	t.Run("Update Article", func(t *testing.T) {
		article := map[string]string{"title": "Updated Article", "body": "This is an updated article."}
		jsonArticle, _ := json.Marshal(article)

		req, _ := http.NewRequest("PUT", "/api/articles/1", bytes.NewBuffer(jsonArticle))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test article deletion
	t.Run("Delete Article", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/articles/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestArticleInteractions(t *testing.T) {
	r := setupRouter()

	// Test like an article
	t.Run("Like Article", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/articles/1/like", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test unlike an article
	t.Run("Unlike Article", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/articles/1/like", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	// Test adding a comment
	t.Run("Add Comment", func(t *testing.T) {
		comment := map[string]string{"body": "This is a test comment."}
		jsonComment, _ := json.Marshal(comment)

		req, _ := http.NewRequest("POST", "/api/articles/1/comments", bytes.NewBuffer(jsonComment))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test retrieving comments
	t.Run("Get Comments", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/articles/1/comments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test deleting a comment
	t.Run("Delete Comment", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/articles/1/comments/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}