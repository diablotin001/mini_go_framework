package user

import (
    "errors"
    "mini_go/pkg/model"
)

func LoginService(username, password string) (*model.User, error) {
    u, err := GetUserByUsername(username)
    if err != nil {
        return nil, err
    }
    // 简化密码比较
    if u.Password != password {
        return nil, errors.New("invalid credentials")
    }
    return u, nil
}

func RegisterService(username, password, email string) error {
    u := &model.User{Username: username, Password: password, Email: email}
    return CreateUser(u)
}
