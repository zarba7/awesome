package persistence

import (
	"github.com/go-redis/redis"
	"pbRole"
	"src/_infrastructure/json"
	store "src/_infrastructure/store/mysql"
)

type event struct {
	c   *MysqlConsumer
	msg *redis.Message
	ev  []*pbRole.Event
}

func (e *event) Unmarshal() error {
	return json.Unmarshal([]byte(e.msg.Payload), &e.ev)
}

func (e *event) Do() {
	var err = e.c.db.Transaction(func(db *store.MySQLClient) error {
		return nil
	})
	if err != nil {
		panic(err)
	}
}
