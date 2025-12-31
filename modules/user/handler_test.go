package user

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "mini_go/pkg/config"
    "mini_go/pkg/model"
    "github.com/gin-gonic/gin"
)

type mockSvcRepo struct{ u *model.User }

func (m mockSvcRepo) GetByUsername(username string) (*model.User, error) { return m.u, nil }

func TestUserHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)
    userRepo = mockSvcRepo{u: &model.User{ID: 1, Username: "amy", Password: "123456"}}
    config.Conf = &config.Config{}
    config.Conf.JWT.Secret = "test"
    config.Conf.JWT.Expire.D = int64(time.Hour)
    config.Conf.JWT.RefreshExpire.D = int64(24 * time.Hour)
    r := gin.Default()
    r.POST("/login", Login)
    payload := map[string]string{"username": "amy", "password": "123456"}
    body, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != 200 {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    var resp struct {
        Code int `json:"code"`
        Data struct {
            AccessToken string `json:"access_token"`
        } `json:"data"`
    }
    _ = json.Unmarshal(w.Body.Bytes(), &resp)
    if resp.Data.AccessToken == "" {
        t.Fatalf("empty access token")
    }
}
