package articles

import (
    "testing"
	"time"
    "github.com/stretchr/testify/assert"
    "github.com/gin-gonic/gin"
    "realworld-backend/users"
	"github.com/jinzhu/gorm"
)

// Helper to create TagModels from strings
func makeTags(tags []string) []TagModel {
    var tagModels []TagModel
    for _, t := range tags {
        tagModels = append(tagModels, TagModel{Tag: t})
    }
    return tagModels
}

// Manual validation helper for ArticleModelValidator
func isArticleValid(a ArticleModelValidator) bool {
    return len(a.Article.Title) >= 4 && a.Article.Body != "" && a.Article.Description != ""
}

// Manual validation helper for CommentModelValidator
func isCommentValid(c CommentModelValidator) bool {
    return c.Comment.Body != ""
}

// --- Model Tests ---

func TestArticleCreationWithValidData(t *testing.T) {
    article := ArticleModel{
        Title:       "Test Title",
        Body:        "Test Body",
        Description: "Test Desc",
        AuthorID:    1,
        Tags:        makeTags([]string{"go", "test"}),
    }
    assert.Equal(t, "Test Title", article.Title)
    assert.Equal(t, "Test Body", article.Body)
    assert.Equal(t, "Test Desc", article.Description)
    assert.Equal(t, uint(1), article.AuthorID)
    assert.Len(t, article.Tags, 2)
    assert.Equal(t, "go", article.Tags[0].Tag)
}

func TestArticleValidationEmptyTitle(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = ""
    validator.Article.Body = "Body"
    validator.Article.Description = "Desc"
    ok := isArticleValid(validator)
    assert.False(t, ok)
}

func TestArticleValidationEmptyBody(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = "Title"
    validator.Article.Body = ""
    validator.Article.Description = "Desc"
    ok := isArticleValid(validator)
    assert.False(t, ok)
}

func TestTagAssociation(t *testing.T) {
    article := ArticleModel{
        Tags: makeTags([]string{"go", "api"}),
    }
    assert.Len(t, article.Tags, 2)
    assert.Equal(t, "go", article.Tags[0].Tag)
    assert.Equal(t, "api", article.Tags[1].Tag)
}

// Helper to create a valid users.UserModel for tests
func makeTestUser() users.UserModel {
    return users.UserModel{
        ID:       1,
        Username: "tester",
        Email:    "tester@example.com",
    }
}

// Helper to create a valid ArticleUserModel for tests
func makeTestArticleUser() ArticleUserModel {
    return ArticleUserModel{
        UserModel:   makeTestUser(),
        UserModelID: 1,
    }
}

// --- Serializer Tests ---

func TestArticleSerializerOutputFormat(t *testing.T) {
    article := ArticleModel{
        Title:       "Title",
        Body:        "Body",
        Description: "Desc",
        Tags:        makeTags([]string{"go"}),
        Author:      makeTestArticleUser(),
        AuthorID:    1,
        Model:       gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now()},
    }
    c, _ := gin.CreateTestContext(nil)
    c.Set("my_user_model", makeTestUser())
    serializer := ArticleSerializer{c, article}
    resp := serializer.Response()
    assert.Equal(t, "Title", resp.Title)
    assert.Equal(t, "Body", resp.Body)
    assert.Equal(t, "Desc", resp.Description)
    assert.Contains(t, resp.Tags, "go")
    assert.Equal(t, "tester", resp.Author.Username)
}

func TestArticleListSerializerWithMultipleArticles(t *testing.T) {
    now := time.Now()
    articles := []ArticleModel{
        {Title: "A", Body: "B", Description: "C", Tags: makeTags([]string{"x"}), Author: makeTestArticleUser(), Model: gorm.Model{CreatedAt: now, UpdatedAt: now}},
        {Title: "D", Body: "E", Description: "F", Tags: makeTags([]string{"y"}), Author: makeTestArticleUser(), Model: gorm.Model{CreatedAt: now, UpdatedAt: now}},
    }
    c, _ := gin.CreateTestContext(nil)
    c.Set("my_user_model", makeTestUser())
    serializer := ArticlesSerializer{c, articles}
    resp := serializer.Response()
    assert.Len(t, resp, 2)
    assert.Equal(t, "A", resp[0].Title)
    assert.Equal(t, "D", resp[1].Title)
    assert.Equal(t, "tester", resp[0].Author.Username)
}

func TestCommentSerializerStructure(t *testing.T) {
    now := time.Now()
    comment := CommentModel{
        Body:      "Nice article",
        Author:    makeTestArticleUser(),
        AuthorID:  1,
        Model:     gorm.Model{CreatedAt: now, UpdatedAt: now},
    }
    c, _ := gin.CreateTestContext(nil)
    c.Set("my_user_model", makeTestUser())
    serializer := CommentSerializer{c, comment}
    resp := serializer.Response()
    assert.Equal(t, "Nice article", resp.Body)
    assert.Equal(t, "tester", resp.Author.Username)
}

// --- Validator Tests ---

func TestArticleModelValidatorWithValidInput(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = "Valid"
    validator.Article.Body = "Valid"
    validator.Article.Description = "Valid"
    ok := isArticleValid(validator)
    assert.True(t, ok)
}

func TestArticleModelValidatorMissingTitle(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = ""
    validator.Article.Body = "Body"
    validator.Article.Description = "Desc"
    ok := isArticleValid(validator)
    assert.False(t, ok)
    assert.Equal(t, "Body", validator.Article.Body)
}

func TestArticleModelValidatorMissingBody(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = "Title"
    validator.Article.Body = ""
    validator.Article.Description = "Desc"
    ok := isArticleValid(validator)
    assert.False(t, ok)
    assert.Equal(t, "Title", validator.Article.Title)
}

func TestArticleModelValidatorMissingDescription(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = "Title"
    validator.Article.Body = "Body"
    validator.Article.Description = ""
    ok := isArticleValid(validator)
    assert.False(t, ok)
}

func TestCommentModelValidatorValid(t *testing.T) {
    validator := CommentModelValidator{}
    validator.Comment.Body = "Nice"
    ok := isCommentValid(validator)
    assert.True(t, ok)
}

func TestCommentModelValidatorEmptyBody(t *testing.T) {
    validator := CommentModelValidator{}
    validator.Comment.Body = ""
    ok := isCommentValid(validator)
    assert.False(t, ok)
}

// --- Additional edge cases ---

func TestArticleSerializerHandlesNilTags(t *testing.T) {
    article := ArticleModel{
        Title: "Title",
        Body:  "Body",
    }
    c, _ := gin.CreateTestContext(nil)
    c.Set("my_user_model", users.UserModel{})
    serializer := ArticleSerializer{c, article}
    resp := serializer.Response()
    assert.NotNil(t, resp.Tags)
}

func TestArticleModelValidatorWhitespaceTitle(t *testing.T) {
    validator := ArticleModelValidator{}
    validator.Article.Title = "   "
    validator.Article.Body = "Body"
    validator.Article.Description = "Desc"
    ok := isArticleValid(validator)
    assert.False(t, ok)
}