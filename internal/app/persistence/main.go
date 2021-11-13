package persistence

import (
	"src/_infrastructure/cache"
	"src/persistence/role"
)

func Main() {
	sub := cache.NewRedisPubSub("127.0.0.1", "1", 0)
	var sub2 = &dao.SubscribeRedis{PubSub: sub, Topic: "topic"}
	sub2.Produce()
	s := &dao.StoreMysql{MysqlConsumer: role.Events}
	s.Consume(nil)
}
