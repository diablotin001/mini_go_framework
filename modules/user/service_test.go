package user

import (
	"errors"
	"mini_go/pkg/config"
	"mini_go/pkg/model"
	"testing"
	"time"
)

type mockRepo struct {
	u   *model.User
	err error
}

func (m mockRepo) GetByUsername(username string) (*model.User, error) {
	return m.u, m.err
}

func TestLoginService_Success(t *testing.T) {
	userRepo = mockRepo{u: &model.User{ID: 1, Username: "amy", Password: "123456"}}
	config.Conf = &config.Config{}
	config.Conf.JWT.Secret = "test"
	config.Conf.JWT.Expire.D = int64(time.Hour)
	config.Conf.JWT.RefreshExpire.D = int64(24 * time.Hour)
	pair, err := LoginService("amy", "123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatalf("empty tokens")
	}
}

func TestLoginService_Fail(t *testing.T) {
	userRepo = mockRepo{err: errors.New("not found")}
	config.Conf = &config.Config{}
	config.Conf.JWT.Secret = "test"
	config.Conf.JWT.Expire.D = int64(time.Hour)
	config.Conf.JWT.RefreshExpire.D = int64(24 * time.Hour)
	_, err := LoginService("x", "y")
	if err == nil {
		t.Fatalf("expected error")
	}
}
