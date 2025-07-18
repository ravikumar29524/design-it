package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CacheItem holds the value and its expiration time
type CacheItem struct {
	value      int
	expiryTime time.Time
}

type Cache struct {
	mu    sync.Mutex
	items map[int]CacheItem
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[int]CacheItem),
		ttl:   ttl,
	}
}

// add a new item to the cache
func (c *Cache) AddItem(value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[value] = CacheItem{
		value:      value,
		expiryTime: time.Now().Add(c.ttl),
	}
}

// check for hit or miss in the cache
func (c *Cache) CheckNumber(value int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.items[value]
	if !exists {
		return false
	}
	if time.Now().After(item.expiryTime) {
		delete(c.items, value) // passive expiry cleanup
		return false
	}
	return true
}

// ActiveExpiry removes expired items regularly
func (c *Cache) ActiveExpiry(interval time.Duration, maxScanDuration time.Duration) {
	for {
		time.Sleep(interval)

		start := time.Now()
		c.mu.Lock()

		for k, v := range c.items {
			if time.Since(start) > maxScanDuration {
				break
				// stop early to reduce lock contention
				// Anyway the passive check is present which will help to
				// clean the expired items at fetch time
			}
			if time.Now().After(v.expiryTime) {
				fmt.Printf("Removing the expired item: %d\n", k)
				delete(c.items, k)
			}
		}

		c.mu.Unlock()
	}
}

// generate random numbers and adds them to the cache
func Producer(c *Cache) {
	for {
		num := rand.Intn(100) + 1 // 1 to 100
		fmt.Printf("[Producer] Generated number: %d\n", num)
		c.AddItem(num)

		// Sleep for a second before generating the next number
		time.Sleep(1 * time.Second)
	}
}

// Consumer queries for random numbers in cache
func Consumer(c *Cache, x int) {
	for {
		for i := 0; i < x; i++ {
			num := rand.Intn(100) + 1
			if c.CheckNumber(num) {
				fmt.Printf("[Consumer] HIT: Found number: %d\n", num)
			} else {
				fmt.Printf("[Consumer] MISS: Number NOT found: %d\n", num)
			}
			time.Sleep(1 * time.Second)
		}
		fmt.Println("[Consumer] Sleeping for 30 seconds before next iteration...")
		time.Sleep(30 * time.Second)
	}
}

func main() {
	// the parameter is the time to live for cache items
	// increase it to see the difference in cache hits
	cache := NewCache(5 * time.Second)

	// Better to clean it actively without impacting the performance
	// 1st parameter (2 seconds ) is the interval to check for expired items
	// 2nd parameter (10 ms) is the maximum time to spend in removing expired items
	// Since cleaning can take time and it holds the lock, so let's keep it short
	go cache.ActiveExpiry(2*time.Second, 10*time.Millisecond)

	// Start producer
	go Producer(cache)

	// 5 -> number of queries per iteration
	go Consumer(cache, 10)

	// keep the main function running
	select {}
}
