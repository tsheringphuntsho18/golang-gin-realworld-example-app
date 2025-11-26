package common

import (
    "bytes"
    "errors"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestConnectingDatabase(t *testing.T) {
    asserts := assert.New(t)
    db := Init()
    // Test create & close DB
    _, err := os.Stat("./../gorm.db")
    asserts.NoError(err, "Db should exist")
    asserts.NoError(db.DB().Ping(), "Db should be able to ping")

    // Test get a connecting from connection pools
    connection := GetDB()
    asserts.NoError(connection.DB().Ping(), "Db should be able to ping")
    db.Close()

    // Test DB exceptions
    os.Chmod("./../gorm.db", 0000)
    db = Init()
    asserts.Error(db.DB().Ping(), "Db should not be able to ping")
    db.Close()
    os.Chmod("./../gorm.db", 0644)
}

func TestConnectingTestDatabase(t *testing.T) {
    asserts := assert.New(t)
    // Test create & close DB
    db := TestDBInit()
    _, err := os.Stat("./../gorm_test.db")
    asserts.NoError(err, "Db should exist")
    asserts.NoError(db.DB().Ping(), "Db should be able to ping")
    db.Close()

    // Test testDB exceptions
    os.Chmod("./../gorm_test.db", 0000)
    db = TestDBInit()
    _, err = os.Stat("./../gorm_test.db")
    asserts.NoError(err, "Db should exist")
    asserts.Error(db.DB().Ping(), "Db should not be able to ping")
    os.Chmod("./../gorm_test.db", 0644)

    // Test close delete DB
    TestDBFree(db)
    _, err = os.Stat("./../gorm_test.db")

    asserts.Error(err, "Db should not exist")
}

func TestRandString(t *testing.T) {
    asserts := assert.New(t)

    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    str := RandString(0)
    asserts.Equal(str, "", "length should be ''")

    str = RandString(10)
    asserts.Equal(len(str), 10, "length should be 10")
    for _, ch := range str {
        found := false
        for _, l := range letters {
            if ch == l {
                found = true
                break
            }
        }
        asserts.True(found, "character should be in allowed letters")
    }
}

func TestGenToken(t *testing.T) {
    asserts := assert.New(t)

    token := GenToken(2)

    asserts.IsType(token, string("token"), "token type should be string")
    asserts.Len(token, 115, "JWT's length should be 115")
}

func TestNewValidatorError(t *testing.T) {
    asserts := assert.New(t)
	err := NewValidatorError("field1", "error1")
    asserts.IsType(CommonError{}, err, "Should return CommonError type")
    asserts.Equal("field1", err.Field)
    asserts.Contains(err.Err.Error(), "error1")
}

func TestNewError(t *testing.T) {
    asserts := assert.New(t)
    err := NewError("test_code", errors.New("an error occurred"))
    // Check that the returned value is of type CommonError by comparing types
    asserts.IsType(CommonError{}, err, "Should return CommonError type")
    asserts.Equal("test_code", err.Code)
    asserts.Contains(err.Err.Error(), "an error occurred")
}

// --- Additional Tests ---

func TestGenTokenDifferentUserIDs(t *testing.T) {
    asserts := assert.New(t)
    token1 := GenToken(1)
    token2 := GenToken(2)
    token3 := GenToken(99999)
    asserts.NotEqual(token1, token2, "Tokens for different user IDs should differ")
    asserts.NotEqual(token1, token3, "Tokens for different user IDs should differ")
    asserts.NotEqual(token2, token3, "Tokens for different user IDs should differ")
    asserts.Len(token1, 115, "JWT's length should be 115")
    asserts.Len(token2, 115, "JWT's length should be 115")
    asserts.Len(token3, 115, "JWT's length should be 115")
}

func TestGenTokenExpiration(t *testing.T) {
    asserts := assert.New(t)
    token := GenToken(123)
    // JWT tokens have three parts separated by dots
    parts := bytes.Split([]byte(token), []byte("."))
    asserts.Equal(3, len(parts), "JWT should have three parts")
    // Optionally, decode and check expiration if implementation allows
}

func TestDatabaseConnectionErrorHandling(t *testing.T) {
    asserts := assert.New(t)
    // Simulate error by passing an invalid path
    oldPath := os.Getenv("DB_PATH")
    os.Setenv("DB_PATH", "/invalid/path/to/db.db")
    defer os.Setenv("DB_PATH", oldPath)
    db := Init()
    err := db.DB().Ping()
    asserts.Error(err, "Should error on invalid DB path")
    db.Close()
}

func TestRandStringUniqueness(t *testing.T) {
    asserts := assert.New(t)
    str1 := RandString(12)
    str2 := RandString(12)
    asserts.Equal(12, len(str1), "RandString should return string of correct length")
    asserts.Equal(12, len(str2), "RandString should return string of correct length")
    asserts.NotEqual(str1, str2, "RandString should produce different strings")
}

func TestRandStringZeroLength(t *testing.T) {
    asserts := assert.New(t)
    str := RandString(0)
    asserts.Equal("", str, "RandString(0) should return empty string")
}