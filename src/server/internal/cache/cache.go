package cache

import (
	"fmt"
	"time"
)

type Cache[T any] struct {
	data        T
	MaxAge      time.Duration
	lastRefresh time.Time
	refresher   func(*T) error
}

func New[T any](ma time.Duration, r func(*T) error) (*Cache[T], error) {
	c := &Cache[T]{
		MaxAge:      ma,
		lastRefresh: time.Now(),
		refresher:   r,
	}
	err := r(&c.data)
	if err != nil {
		return nil, fmt.Errorf("Failed to pull data with given refresher function, %s", err)
	}
	return c, nil
}

func (c *Cache[T]) GetData() (T, error) {
	if time.Since(c.lastRefresh) >= c.MaxAge {
		err := c.refresher(&c.data)
		if err != nil {
			return c.data, err
		}
		c.lastRefresh = time.Now()
	}
	return c.data, nil
}

func (c *Cache[T]) MustGetData() T {
	if time.Since(c.lastRefresh) >= c.MaxAge {
		if err := c.refresher(&c.data); err != nil {
			panic(err)
		}
		c.lastRefresh = time.Now()
	}
	return c.data
}

func (c *Cache[T]) SetData(d T) {
	c.data = d
	c.lastRefresh = time.Now()
}

func (c *Cache[T]) GetLastRefresh() time.Time {
	return c.lastRefresh
}

func (c *Cache[T]) Refresh() error {
	if err := c.refresher(&c.data); err != nil {
		return err
	}
	c.lastRefresh = time.Now()
	return nil
}
