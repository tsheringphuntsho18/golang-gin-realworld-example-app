package articles

import (
	"errors"
	"regexp"
)

// Article represents an article structure
type Article struct {
	Title   string
	Content string
	Author  string
}

// ValidateArticle checks if the article fields are valid
func ValidateArticle(article Article) error {
	if err := validateTitle(article.Title); err != nil {
		return err
	}
	if err := validateContent(article.Content); err != nil {
		return err
	}
	if err := validateAuthor(article.Author); err != nil {
		return err
	}
	return nil
}

func validateTitle(title string) error {
	if len(title) == 0 {
		return errors.New("title cannot be empty")
	}
	if len(title) > 100 {
		return errors.New("title cannot exceed 100 characters")
	}
	return nil
}

func validateContent(content string) error {
	if len(content) == 0 {
		return errors.New("content cannot be empty")
	}
	if len(content) > 1000 {
		return errors.New("content cannot exceed 1000 characters")
	}
	return nil
}

func validateAuthor(author string) error {
	if len(author) == 0 {
		return errors.New("author cannot be empty")
	}
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", author); !matched {
		return errors.New("author can only contain alphanumeric characters and underscores")
	}
	return nil
}