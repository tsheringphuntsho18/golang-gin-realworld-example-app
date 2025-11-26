package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"realworld-backend/articles"
	"realworld-backend/common"
	"realworld-backend/users"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := common.InitTestDB() // You may need to implement this for test DB isolation
	common.DB = db
	r := gin.Default()
	users.Routers(r)
	articles.Routers(r)
	return r
}

func TestUserRegistration(t *testing.T) {
	router := setupRouter()
	payload := `{"user":{"username":"testuser","email":"testuser@example.com","password":"password"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["user"])
	assert.NotEmpty(t, resp["user"].(map[string]interface{})["token"])
}

func TestUserRegistrationSavedInDB(t *testing.T) {
	router := setupRouter()
	payload := `{"user":{"username":"dbuser","email":"dbuser@example.com","password":"password"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var user users.UserModel
	common.DB.Where("email = ?", "dbuser@example.com").First(&user)
	assert.Equal(t, "dbuser", user.Username)
}

func TestUserLoginSuccess(t *testing.T) {
	router := setupRouter()
	// Register first
	payload := `{"user":{"username":"loginuser","email":"loginuser@example.com","password":"password"}}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)
	// Login
	loginPayload := `{"user":{"email":"loginuser@example.com","password":"password"}}`
	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["user"].(map[string]interface{})["token"])
}

func TestUserLoginInvalidCredentials(t *testing.T) {
	router := setupRouter()
	loginPayload := `{"user":{"email":"notfound@example.com","password":"wrong"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 422, w.Code)
}

func TestGetCurrentUserWithValidToken(t *testing.T) {
	router := setupRouter()
	// Register and login
	payload := `{"user":{"username":"meuser","email":"meuser@example.com","password":"password"}}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["user"].(map[string]interface{})["token"].(string)
	// Get current user
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/user", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetCurrentUserWithInvalidToken(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user", nil)
	req.Header.Set("Authorization", "Token invalidtoken")
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestGetCurrentUserWithoutToken(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func registerAndLogin(router *gin.Engine, username, email, password string) string {
	payload := fmt.Sprintf(`{"user":{"username":"%s","email":"%s","password":"%s"}}`, username, email, password)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["user"].(map[string]interface{})["token"].(string)
}

func TestCreateArticleAuthenticated(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author1", "author1@example.com", "password")
	articlePayload := `{"article":{"title":"Test Article","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Test Article", resp["article"].(map[string]interface{})["title"])
}

func TestCreateArticleWithoutAuth(t *testing.T) {
	router := setupRouter()
	articlePayload := `{"article":{"title":"No Auth","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestListArticles(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotNil(t, resp["articles"])
}

func TestGetSingleArticle(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author2", "author2@example.com", "password")
	articlePayload := `{"article":{"title":"Single Article","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Get article
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/articles/"+slug, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUpdateArticleByAuthor(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author3", "author3@example.com", "password")
	articlePayload := `{"article":{"title":"Update Me","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Update
	updatePayload := `{"article":{"title":"Updated Title"}}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/articles/"+slug, bytes.NewBufferString(updatePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestUpdateArticleUnauthorized(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author4", "author4@example.com", "password")
	articlePayload := `{"article":{"title":"No Update","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Try update with no token
	updatePayload := `{"article":{"title":"Should Not Update"}}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/articles/"+slug, bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestDeleteArticleByAuthor(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author5", "author5@example.com", "password")
	articlePayload := `{"article":{"title":"Delete Me","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/articles/"+slug, nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestDeleteArticleUnauthorized(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "author6", "author6@example.com", "password")
	articlePayload := `{"article":{"title":"No Delete","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Try delete with no token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/articles/"+slug, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func TestFavoriteUnfavoriteArticle(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "favuser", "favuser@example.com", "password")
	articlePayload := `{"article":{"title":"Fav Article","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Favorite
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/articles/"+slug+"/favorite", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Unfavorite
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/articles/"+slug+"/favorite", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestArticleCommentsCRUD(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(router, "commenter", "commenter@example.com", "password")
	articlePayload := `{"article":{"title":"Comment Article","description":"desc","body":"body"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(articlePayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	slug := resp["article"].(map[string]interface{})["slug"].(string)
	// Add comment
	commentPayload := `{"comment":{"body":"Nice article!"}}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/articles/"+slug+"/comments", bytes.NewBufferString(commentPayload))
	req.Header.Set("Authorization", "Token "+token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var commentResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &commentResp)
	commentID := int(commentResp["comment"].(map[string]interface{})["id"].(float64))
	// List comments
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/articles/"+slug+"/comments", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Delete comment
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/articles/%s/comments/%d", slug, commentID), nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

// Clean up test DB if needed
func TestMain(m *testing.M) {
	code := m.Run()
	os.Remove("test.db")
	os.Exit(code)
}