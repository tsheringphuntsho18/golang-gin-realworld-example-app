package articles

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
	// Assume we have a function to register routes
	// registerRoutes(r)
	return r
}

func TestCreateArticle(t *testing.T) {
	router := setupRouter()

	article := map[string]interface{}{
		"title":   "Test Article",
		"content": "This is a test article.",
	}

	jsonValue, _ := json.Marshal(article)
	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetArticle(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateArticle(t *testing.T) {
	router := setupRouter()

	article := map[string]interface{}{
		"title":   "Updated Article",
		"content": "This is an updated test article.",
	}

	jsonValue, _ := json.Marshal(article)
	req, _ := http.NewRequest("PUT", "/articles/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteArticle(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/articles/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestListArticles(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/articles", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserAuthentication(t *testing.T) {
	router := setupRouter()

	user := map[string]interface{}{
		"username": "testuser",
		"password": "password",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserRegistration(t *testing.T) {
	router := setupRouter()

	user := map[string]interface{}{
		"username": "newuser",
		"password": "newpassword",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUserProfile(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/users/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUserProfile(t *testing.T) {
	router := setupRouter()

	user := map[string]interface{}{
		"username": "updateduser",
		"email":    "updated@example.com",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/profile", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleInteractionLike(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/articles/1/like", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestArticleInteractionComment(t *testing.T) {
	router := setupRouter()

	comment := map[string]interface{}{
		"body": "This is a comment.",
	}

	jsonValue, _ := json.Marshal(comment)
	req, _ := http.NewRequest("POST", "/articles/1/comments", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetArticleComments(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/articles/1/comments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteComment(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/articles/1/comments/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetArticleBySlug(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/articles/test-article", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetArticlesByTag(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/articles?tag=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}