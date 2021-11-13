package domain

import (
	"ddd/adaptor"
	"log"
	"sync"
	"time"
)

type cache struct {
	mu        sync.Mutex
	hash      map[int64]*aggregate
	pool      sync.Pool
	chRecycle chan *aggregate
}

func (the *cache) run(repo adaptor.DomainRepo) {
	the.chRecycle = make(chan *aggregate, 8)
	the.pool.New = func() interface{} { return &aggregate{} }
	the.hash = make(map[int64]*aggregate)

	go func() {
		var tick = time.Tick(59 * time.Second)
		for range tick {
			the.clear()
		}
	}()

	var err error
	for agg := range the.chRecycle {
		if agg.po.SnapshotVersion < agg.Version() {
			if agg.po.Snapshot, err = agg.ss.Save(); err != nil {
				log.Printf("[%s %d] recycle Snapshot %v", agg.Tag(), agg.Root(), err)
			} else if err = repo.UpdateSnapshot(&agg.po); err != nil {
				log.Printf("[%s %d] recycle Snapshot %v", agg.Tag(), agg.Root(), err)
			}
		}
		the.recycle(agg)
	}
}

func (the *cache) clear() {
	the.mu.Lock()
	defer the.mu.Unlock()
	for _, agg := range the.hash {
		if isExpired(agg) {
			delete(the.hash, agg.Root())
			the.chRecycle <- agg
		}
	}
}

func (the *cache) Remove(agg *aggregate) {
	the.mu.Lock()
	defer the.mu.Unlock()
	delete(the.hash, agg.Root())
	the.recycle(agg)
}

func (the *cache) recycle(agg *aggregate) {
	agg.clear()
	agg.expired = 0
	the.pool.Put(agg)
}

func isExpired(agg *aggregate) bool {
	return 0 < agg.expired && agg.expired < time.Now().Unix()
}

func (the *cache) Find(root int64) (agg *aggregate, ok bool) {
	the.mu.Lock()
	defer the.mu.Unlock()
	agg, ok = the.hash[root]
	if !ok {
		agg = the.pool.Get().(*aggregate)
		the.hash[root] = agg
	}
	agg.expired = time.Now().Add(time.Hour / 2).Unix()
	return
}
