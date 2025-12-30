package user

import (
	"errors"
	"mini_go/pkg/cache"
	"mini_go/pkg/config"
	"mini_go/pkg/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func LoginService(username, password string) (*TokenPair, error) {
	u, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	// fmt.Println(u)
	if u.Password != password {
		return nil, errors.New("invalid credentials")
	}
	access, err := createToken(u.ID, "access", time.Duration(config.Conf.JWT.Expire.D))
	if err != nil {
		return nil, err
	}
	refresh, err := createToken(u.ID, "refresh", time.Duration(config.Conf.JWT.RefreshExpire.D))
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func RegisterService(username, password, email string) error {
	u := &model.User{Username: username, Password: password, Email: email}
	return CreateUser(u)
}

func createToken(uid uint, typ string, ttl time.Duration) (string, error) {
	jti := uuid.NewString()
	exp := time.Now().Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"typ": typ,
		"jti": jti,
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(config.Conf.JWT.Secret))
}

func LogoutService(accessToken string, refreshToken string) error {
	t, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.Conf.JWT.Secret), nil })
	if err == nil && t.Valid {
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			if jti, ok2 := claims["jti"].(string); ok2 {
				var exp int64
				if ev, ok3 := claims["exp"].(float64); ok3 {
					exp = int64(ev)
				}
				// TODO: 基于uid封锁
				// uid: 1 map[exp:1.76733921e+09 iat:1.76708001e+09 jti:43a393e8-fe2a-4544-b33e-87d91001c10a typ:access uid:1]
				// fmt.Println("uid:", claims["uid"], claims)
				ttl := time.Until(time.Unix(exp, 0))
				if ttl > 0 {
					_ = cache.SetString("jwt:blacklist:"+jti, "1", ttl)
				}
			}
		}
	}
	if refreshToken != "" {
		t2, err2 := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.Conf.JWT.Secret), nil })
		if err2 == nil && t2.Valid {
			if claims, ok := t2.Claims.(jwt.MapClaims); ok {
				if jti, ok2 := claims["jti"].(string); ok2 {
					var exp int64
					if ev, ok3 := claims["exp"].(float64); ok3 {
						exp = int64(ev)
					}
					ttl := time.Until(time.Unix(exp, 0))
					if ttl > 0 {
						_ = cache.SetString("jwt:blacklist:"+jti, "1", ttl)
					}
				}
			}
		}
	}
	return nil
}

func RefreshService(refreshToken string) (string, error) {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.Conf.JWT.Secret), nil })
	if err != nil || !t.Valid {
		return "", errors.New("invalid refresh token")
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid refresh token")
	}
	if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
		return "", errors.New("invalid token type")
	}
	var uid uint
	switch v := claims["uid"].(type) {
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	}
	return createToken(uid, "access", time.Duration(config.Conf.JWT.Expire.D))
}
