//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	ch     chan struct{}
	chStop chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	ch := make(chan struct{})
	chStop := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case ch <- struct{}{}:
				select {
				case <-time.After(interval / time.Duration(maxCount)):
				case <-chStop:
					return
				}
			case <-chStop:
				return
			}
		}
	}()

	return &Limiter{ch, chStop}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case _, ok := <-l.ch:
		if ok {
			return nil
		} else {
			return ErrStopped
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() {
	l.chStop <- struct{}{}
}
