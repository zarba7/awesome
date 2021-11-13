package log

import (
	"github.com/go-redis/redis"
	"io"
	"log"
	"time"
)

// 为 logger 提供写入 redis 队列的 io 接口
type redisWriter struct {
	ListKey string
	*redis.Client
	MaxLen int64
}

func (w *redisWriter) Write(p []byte) (int, error) {
	n, err := w.Client.LPush(w.ListKey, p).Result()
	if err != nil {
		return int(n), err
	}
	if n > w.MaxLen {
		_, err = w.Client.LTrim(w.ListKey, 0, w.MaxLen-1).Result()
	}
	return int(w.MaxLen), err
}

func NewRedisWriter(name string, addr, auth string, db int) io.Writer {
	cli := redis.NewClient(&redis.Options{
		Addr:               addr,
		Password:           auth,
		DB:                 db,
		PoolSize:           8,
		PoolTimeout:        10 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
		OnConnect:          nil,
	})
	if err := cli.Ping().Err(); err != nil {
		log.Panicf("ping redis %v", err)
	}

	return &redisWriter{
		ListKey: "log:" + name,
		Client:  cli,
		MaxLen:  20000,
	}
}
