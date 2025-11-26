package common

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
    var err error
    DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
    if err != nil {
        return err
    }
    return nil
}

func CloseDatabase() error {
    db, err := DB.DB()
    if err != nil {
        return err
    }
    return db.Close()
}