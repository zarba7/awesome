package persistence

import (
	"ddd/adaptor"
	"ddd/po"
	"fmt"
	"github.com/go-redis/redis"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

func NewRedisRepo(addr, auth string, db int) (adaptor.DomainRepo, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:               addr,
		Password:           auth,
		DB:                 db,
		PoolSize:           8,
		PoolTimeout:        10 * time.Second,
		IdleTimeout:        time.Minute,
		IdleCheckFrequency: 100 * time.Millisecond,
		OnConnect: func(conn *redis.Conn) error {
			zap.L().Info("connect to redis", zap.String("addr", conn.String()), zap.String("auth", auth), zap.Int("db", db))
			return conn.Ping().Err()
		},
	})
	repo := &redisRepo{Client: cli}
	err := repo.init()
	return repo, err
}

type redisRepo struct {
	saveTrx *redis.Script
	//updateSnapshotTrx    *redis.Script
	loadAllTrx *redis.Script
	loadTrx    *redis.Script
	snapTrx    *redis.Script
	*redis.Client
}

func (repo *redisRepo) init() error { //list
	repo.loadAllTrx = redis.NewScript(`
		local kA, kE = KEYS[1], KEYS[2]
		local agg = redis.call("HGetAll", kA)
		local allEvents = redis.call("LRange", kE, 0, -1)
		return agg, allEvents
	`)
	err := repo.loadAllTrx.Load(repo.Client).Err()
	repo.loadTrx = redis.NewScript(`
		local kA, kE = KEYS[1], KEYS[2]
		local ver = ARGV[1]
		local newVer = redis.call("HGet", kA, "version")
		local allEvents = redis.call("LRange", kE, newVer-ver, -1)
		return newVer, allEvents
	`)
	err = repo.loadTrx.Load(repo.Client).Err()
	repo.saveTrx = redis.NewScript(`
		local kA, kE, snapVer, events = KEYS[1], KEYS[2], ARGV[1], ARGV[2]
		local newVer = redis.call("HIncrBy", kA, "version", 1)
		if snapVer + 1 ~= newVer then
			redis.call("HIncrBy", kA, "version", -1)
			return -1
		end
		redis.call("RPush", kE, events)
		return 0
	`)
	err = repo.saveTrx.Load(repo.Client).Err()

	/*repo.updateSnapshotTrx = redis.NewScript(`
		local kA, snapVer, snapshot = KEYS[1], ARGV[1], ARGV[2]
		redis.call("HMSet", kA, "snapshotVersion", snapVer, "snapshot", snapshot)
	`)
	err = repo.updateSnapshotTrx.Load(repo.Client).Err()*/
	return err
}

func key2(agg *po.DomainAggregate) []string { return []string{kAgg(agg), kEvt(agg)} }
func kEvt(agg *po.DomainAggregate) string   { return fmt.Sprintf("{%s%d}:lEvent", agg.Tag, agg.ID) }
func kAgg(agg *po.DomainAggregate) string   { return fmt.Sprintf("{%s%d}:hAggregate", agg.Tag, agg.ID) }

func (repo *redisRepo) UpdateSnapshot(agg *po.DomainAggregate) error {
	err := repo.HMSet(kAgg(agg), map[string]interface{}{
		"snapshotVersion": agg.SnapshotVersion,
		"snapshot":        agg.Snapshot,
	}).Err()
	return err
}

func (repo *redisRepo) LoadAll(agg *po.DomainAggregate) error {
	result, err := repo.loadAllTrx.EvalSha(repo.Client, key2(agg)).Result()
	if err != nil {
		return err
	}
	zap.L().Debug("", zap.Any("result", result))
	return nil
}
func (repo *redisRepo) Load(agg *po.DomainAggregate) error {
	return repo.loadTrx.EvalSha(repo.Client, key2(agg), agg.CurrVersion).Err()
}

func (repo *redisRepo) Save(agg *po.DomainAggregate) error {
	events, _ := json.Marshal(agg.Events)
	return repo.saveTrx.EvalSha(repo.Client, key2(agg), agg.SnapshotVersion, events).Err()
}
