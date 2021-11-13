package mutex

type locker struct {
	key string
	fac *factory
	mtx *mutex
}

func (the *locker) Lock() {
	the.mtx, _ = the.fac.lock(the.key)
}
func (the *locker) Unlock() {
	if the.fac.unlock(the.key, the.mtx) {
		the.fac.pool.Put(the.mtx)
	}
	the.fac.lockerPool.Put(the)
}
