package articles

import (
	"time"
)

// Article represents the structure of an article in the application.
type Article struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ArticleResponse represents the response structure for an article.
type ArticleResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    AuthorResponse `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthorResponse represents the response structure for an author.
type AuthorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreateArticleRequest represents the request structure for creating an article.
type CreateArticleRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// UpdateArticleRequest represents the request structure for updating an article.
type UpdateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ArticleListResponse represents the response structure for a list of articles.
type ArticleListResponse struct {
	Articles []ArticleResponse `json:"articles"`
	Count    int               `json:"count"`
}