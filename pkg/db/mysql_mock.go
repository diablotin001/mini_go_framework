package db

import (
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func NewMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("mock db err: %v", err)
    }
    dial := mysql.New(mysql.Config{Conn: mockDB, SkipInitializeWithVersion: true})
    db, err := gorm.Open(dial, &gorm.Config{})
    if err != nil {
        t.Fatalf("gorm open err: %v", err)
    }
    return db, mock
}
