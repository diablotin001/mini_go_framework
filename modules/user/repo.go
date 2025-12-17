package user

import (
    "mini_go/pkg/cache"
    "mini_go/pkg/model"
    "mini_go/internal/database"
    "time"
)

func GetUserByUsername(username string) (*model.User, error) {
    key := "user:username:" + username
    if s, err := cache.GetString(key); err == nil && s != "" {
        return &model.User{Username: s}, nil
    }
    var u model.User
    if err := database.DB.Where("username = ?", username).First(&u).Error; err != nil {
        return nil, err
    }
    _ = cache.SetString(key, u.Username, 3600*time.Second)
    return &u, nil
}

func CreateUser(u *model.User) error {
    return database.DB.Create(u).Error
}
