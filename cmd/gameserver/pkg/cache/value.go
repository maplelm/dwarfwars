package cache

import (
	"sync"
	"time"
)

type Refresher[T any] interface {
	Update(*Value[T]) error
}

type Value[T any] struct {
	refresher Refresher[T]

	valueMut   sync.RWMutex
	data       T
	maxAge     time.Duration
	lastPolled time.Time
}

func (v *Value[T]) SetData(val T) {
	v.valueMut.Lock()
	defer v.valueMut.Unlock()
	v.data = val
	v.lastPolled = time.Now()
}

func (v *Value[T]) GetData() (T, error) {
	var e error
	if time.Since(v.lastPolled) >= v.maxAge {
		v.valueMut.Lock()
		e = v.refresher.Update(v)
		v.lastPolled = time.Now()
		v.valueMut.Unlock()
	}
	v.valueMut.RLock()
	defer v.valueMut.RUnlock()
	return v.data, e
}

func (v *Value[T]) SetMaxAge(t time.Duration) {
	v.valueMut.Lock()
	defer v.valueMut.Unlock()
	v.maxAge = t
}

func (v *Value[T]) IncreaseMaxAge(t time.Duration) {
	v.valueMut.Lock()
	defer v.valueMut.Unlock()
	v.maxAge += t
}

func (v *Value[T]) DecreaseMaxAge(t time.Duration) {
	v.valueMut.Lock()
	defer v.valueMut.Unlock()
	v.maxAge -= t
}
