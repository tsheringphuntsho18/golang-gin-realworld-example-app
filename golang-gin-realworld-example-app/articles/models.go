package articles

import (
	"time"
)

// Article represents an article in the system.
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Comment represents a comment on an article.
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// Like represents a like on an article.
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}