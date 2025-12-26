package user

import (
	"errors"
	"mini_go/pkg/config"
	"mini_go/pkg/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func LoginService(username, password string) (string, error) {
	u, err := GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	// 简化密码比较, 生产环境用 bcrypt
	if u.Password != password {
		return "", errors.New("invalid credentials")
	}
	exp := time.Now().Add(time.Duration(config.Conf.JWT.Expire.D))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": u.ID,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(config.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func RegisterService(username, password, email string) error {
	u := &model.User{Username: username, Password: password, Email: email}
	return CreateUser(u)
}
