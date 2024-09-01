package settings

import (
	"time"
)

type TimeCached[T any] struct {
	pollRate    time.Duration
	lastPoll    time.Time
	data        T
	refreshFunc func(*TimeCached[T]) error
}

func NewTimeCached[T any](pr time.Duration, rf func(*TimeCached[T]) error) *TimeCached[T] {
	return &TimeCached[T]{
		pollRate:    pr,
		lastPoll:    time.Unix(0, 0),
		refreshFunc: rf,
	}
}

func (c *TimeCached[T]) Get() (v T, err error) {
	if time.Since(c.lastPoll) >= c.pollRate {
		if err = c.refreshFunc(c); err != nil {
			return
		}
	}
	v = c.data
	return
}

func (c *TimeCached[T]) SetPollRate(r time.Duration) {
	c.pollRate = r
}
