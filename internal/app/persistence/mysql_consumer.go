package persistence

import (
	"github.com/go-redis/redis"
	"src/_ddd/persistence"
	"src/_infrastructure/cache"
	store "src/_infrastructure/store/mysql"
	"src/persistence/dao"
)

type redisPublisher struct {
	*cache.RedisPubSub
	topic string
}

type MysqlConsumer struct {
	persistence.DaoEvents
	db  *store.MySQLClient
	pub redisPublisher
	events dao.Table
}
func (s *MysqlConsumer) SetEvents(tab dao.Table){

	s.events = tab
}

func (s *MysqlConsumer) Run(db *store.MySQLClient, pub redisPublisher) {
	s.db = db
	s.pub = pub
	s.P = s
	s.C = s
	s.DaoEvents.Run(8)
}

func (s *MysqlConsumer) Produce(c persistence.Consumer) {
	s.pub.Subscribe(func(msg *redis.Message) {
		c.Consume(&event{c: s, msg: msg})
	}, s.pub.topic)
}

func (s *MysqlConsumer) Close() error {
	s.DaoEvents.Close()
	s.pub.Close()
	return s.db.Close()
}
