package excache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(3, 2, 1*time.Second)
	go func() {
		cache.Set("mykey", []string{"hello", "world"})
	}()
	go func() {
		cache.Set("mykey", []string{"hello", "exchange"})
	}()
	go func() {
		if _, ok := cache.Get("mykey"); !ok {
			t.Fail()
		}
		if _, ok := cache.Get("mykey"); !ok {
			t.Fail()
		}
		if _, ok := cache.Get("mykey"); ok {
			t.Fail()
		}
	}()
	time.Sleep(1 * time.Second)
	if _, ok := cache.Get("mykey"); ok {
		t.Fail()
	}
}

func BenchmarkLRUCache(b *testing.B) {
	cache := NewLRUCache(2000, 1, 1*time.Second)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		key := fmt.Sprintf("key_%d", i)
		go func() {
			cache.Set(key, "some value")
			cache.Get(key)
			wg.Done()
		}()
	}
	wg.Wait()
}
