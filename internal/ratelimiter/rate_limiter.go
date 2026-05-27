package ratelimiter

import (
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64 // unix nanoseconds, updated atomically
}

type Ratelimiter interface {
	Allow(ip string) bool
}

type BucketTokenRateLimiter struct {
	sync.RWMutex
	clients map[string]*Client
	limit   rate.Limit
	burst   int
}

func NewBucketTokenRatelimiter(burst, limit int) *BucketTokenRateLimiter {
	rl := &BucketTokenRateLimiter{
		clients: make(map[string]*Client),
		burst:   burst,
		limit:   rate.Limit(limit),
	}
	go rl.cleanupStaleClients(5 * time.Minute)
	return rl
}

func (br *BucketTokenRateLimiter) Allow(ip string) bool {
	return br.getLimiter(ip).Allow()
}

func (br *BucketTokenRateLimiter) getLimiter(ip string) *rate.Limiter {
	// Fast path: check under read lock
	br.RLock()
	client, exists := br.clients[ip]
	br.RUnlock()

	if exists {
		// Safe: atomic update, no write lock needed
		client.lastSeen.Store(time.Now().UnixNano())
		return client.limiter
	}

	// Slow path: create new client under write lock
	br.Lock()
	defer br.Unlock()

	// Double-check after acquiring write lock — another goroutine may have
	// already inserted this IP between our RUnlock and Lock above.
	if client, exists = br.clients[ip]; exists {
		client.lastSeen.Store(time.Now().UnixNano())
		return client.limiter
	}

	newClient := &Client{limiter: rate.NewLimiter(br.limit, br.burst)}
	newClient.lastSeen.Store(time.Now().UnixNano())
	br.clients[ip] = newClient
	return newClient.limiter
}

// cleanupStaleClients removes entries that haven't been seen within ttl.
// Runs as a background goroutine started by NewBucketTokenRatelimiter.
func (br *BucketTokenRateLimiter) cleanupStaleClients(ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-ttl).UnixNano()
		br.Lock()
		for ip, client := range br.clients {
			if client.lastSeen.Load() < cutoff {
				delete(br.clients, ip)
			}
		}
		br.Unlock()
	}
}
