package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(addr, password string, db int) error {
	RDB = redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetString(key string) (string, error) {
	// if RDB == nil {
	// 	return "", nil
	// }
	return RDB.Get(Ctx, key).Result()
}

func SetString(key string, value interface{}, ttl time.Duration) error {
	// if RDB == nil {
	// 	return nil
	// }
	return RDB.Set(Ctx, key, value, ttl).Err()
}
