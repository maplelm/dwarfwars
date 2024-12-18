package cache

import (
	"sync"
	"time"
)

type Cache[T any] struct {
	PollRate time.Duration

	lastPolled time.Time

	mutex   sync.Mutex
	data    T
	refresh func(*T) error
}

func New[T any](pr time.Duration, f func(*T) error) *Cache[T] {
	return &Cache[T]{
		PollRate:   pr,
		refresh:    f,
		lastPolled: time.Unix(0, 0),
	}
}

func (c *Cache[T]) Get() (*T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if time.Since(c.lastPolled) >= c.PollRate {
		err := c.refresh(&c.data)
		if err != nil {
			return nil, err
		}
		c.lastPolled = time.Now()
	}
	return &c.data, nil
}

func (c *Cache[T]) ForceRefresh() (*T, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	err := c.refresh(&c.data)
	if err == nil {
		c.lastPolled = time.Now()
	}
	return &c.data, err
}

func (c *Cache[T]) MustGet() *T {
	v, err := c.Get()
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Cache[T]) Set(v T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = v
	c.lastPolled = time.Now()
}
