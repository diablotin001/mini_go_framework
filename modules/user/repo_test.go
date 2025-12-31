package user

import (
    "regexp"
    "testing"
    "mini_go/internal/database"
    mdb "mini_go/pkg/db"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserByUsername(t *testing.T) {
    db, mock := mdb.NewMock(t)
    database.DB = db
    rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "amy", "123456")
    mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).WithArgs("amy", 1).WillReturnRows(rows)
    u, err := GetUserByUsername("amy")
    if err != nil {
        t.Fatalf("unexpected err: %v", err)
    }
    if u.Username != "amy" {
        t.Fatalf("wrong username: %v", u.Username)
    }
}
