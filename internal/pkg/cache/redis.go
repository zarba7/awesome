package cache

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"time"
)

func NewRedisClient(addr, auth string, db int) *RedisClient {	return &RedisClient{redisCli(addr,auth,db)}}

type RedisClient struct {
	*redis.Client
}

func (the *RedisClient) Pipeline(fun func(pipe redis.Pipeliner)) ([]redis.Cmder, error) {
	pipe := the.Client.Pipeline()
	fun(pipe)
	return pipe.Exec()
}

func redisCli(addr, auth string, db int) *redis.Client {
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
		zap.S().Panicf("ping redis %v", err)
	}
	zap.S().Info("connect to redis", zap.String("addr", addr), zap.String("auth", auth), zap.Int("db", db))
	return cli
}

func NewRedisPubSub(addr, auth string, db int) *RedisPubSub {	return &RedisPubSub{redisCli(addr, auth, db)} }

func NewRedisMultiLocker(cli *redis.Client) *RedisMultiLocker { return &RedisMultiLocker{cli: cli} }


/****************************************************/
/****************************************************/
type RedisPubSub struct {
	cli *redis.Client
}

func (the *RedisPubSub) Close() error {
	return the.cli.Close()
}

func (the *RedisPubSub) Subscribe(fn func(msg *redis.Message), channels ...string) {
	ps := the.cli.Subscribe(channels...)
	go redisPubSubReceive(ps, fn)
}

func redisPubSubReceive(ps *redis.PubSub, fn func(msg *redis.Message)) {
	for {
		msg, err := ps.ReceiveMessage()
		if err != nil {
			zap.S().Info(err)
			break
		}
		fn(msg)
	}
}

func (the *RedisPubSub) Publish(channel string, data []byte) error {
	result := the.cli.Publish(channel, data)
	return result.Err()
}

/****************************************************/
/****************************************************/
type RedisMultiLocker struct {
	Prefix string
	cli    *redis.Client
}

func (the *RedisMultiLocker) Lock(ttl time.Duration, keys ...string) bool {
	var n = len(keys)
	if n == 0 {
		return true
	}
	var pairs = make([]interface{}, 0, n*2)
	for i := range keys {
		keys[i] = the.Prefix + keys[i]
		pairs = append(pairs, keys[i], "1")
	}
	var ok, err = the.cli.MSetNX(pairs...).Result()
	if err != nil {
		zap.S().Error(err)
		return false
	}
	if !ok {
		return false
	}
	if ttl < time.Second {
		ttl = time.Second
	}
	for _, k := range keys {
		_, err = the.cli.Expire(k, ttl).Result()
		if err != nil {
			zap.S().Error(err)
			the.cli.Del(keys...)
			return false
		}
	}
	return true
}

func (the *RedisMultiLocker) Unlock(keys ...string) {
	for i := range keys {
		keys[i] = the.Prefix + keys[i]
	}
	var _, err = the.cli.Del(keys...).Result()
	if err != nil {
		zap.S().Error(err)
	}
}
