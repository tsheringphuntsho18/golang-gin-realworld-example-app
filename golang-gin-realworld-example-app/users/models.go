package users

import (
	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
}

// TableName overrides the table name used by User to `users`.
func (User) TableName() string {
	return "users"
}