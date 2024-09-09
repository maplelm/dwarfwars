package settings

import (
	"fmt"
	"time"
)

type CachedValue[T any] struct {
	data        *T
	PollRate    time.Duration
	lastpolled  time.Time
	refreshFunc func(c *CachedValue[T]) error
}

func NewCachedValue[T any](pr time.Duration, rf func(c *CachedValue[T]) error) *CachedValue[T] {
	return &CachedValue[T]{
		data:        nil,
		PollRate:    pr,
		lastpolled:  time.Unix(0, 0),
		refreshFunc: rf,
	}
}

func (c *CachedValue[T]) Get() (*T, error) {
	if time.Since(c.lastpolled) >= c.PollRate {
		err := c.refreshFunc(c)
		if err != nil {
			return c.data, fmt.Errorf("Refresh Failed, potentially stale data, %s", err)
		}
		c.lastpolled = time.Now()
	}
	return c.data, nil

}

func (c *CachedValue[T]) Set(value *T) {
	c.data = value
}
