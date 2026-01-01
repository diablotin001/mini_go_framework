package cache

import (
	"context"
	"fmt"
	"mini_go/logger"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis(addr, password string, db int) error {
	RDB = redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	RDB.AddHook(&logHook{})
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

type logHook struct{}

func (h *logHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		start := time.Now()
		conn, err := next(ctx, network, addr)
		logger.Log.Info("RedisDial", zap.String("addr", addr), zap.Duration("cost", time.Since(start)))
		return conn, err
	}
}

func (h *logHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmd)
		logger.Log.Info("Redis", zap.String("cmd", cmd.FullName()), zap.String("args", fmt.Sprint(cmd.Args())), zap.Duration("cost", time.Since(start)))
		return err
	}
}

func (h *logHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()
		err := next(ctx, cmds)
		logger.Log.Info("RedisPipeline", zap.Int("count", len(cmds)), zap.Duration("cost", time.Since(start)))
		return err
	}
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
