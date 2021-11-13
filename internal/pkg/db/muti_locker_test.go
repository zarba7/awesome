package db

import (
	"sync"
	"testing"
	"time"
)

func TestNewMultiLocker(t *testing.T) {
	var cli = New("127.0.0.1:6379", "", 4)
	var lock = NewMultiLocker(cli.Client)
	var n int
	var lst []int
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			for !lock.Lock(time.Second, "a", "b") {
			}
			lst = append(lst, n)
			n += 1
			lock.Unlock("a", "b")
			wg.Done()
		}()
	}
	wg.Wait()
	for k, v := range lst {
		if k != v {
			println(k, v)
		}
	}
}
