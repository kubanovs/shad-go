//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

type Cond struct {
	L       Locker
	mu      chan struct{}   // внутренний мьютекс для защиты списка
	waiters []chan struct{} // по одному каналу на каждый ожидающий горутин
}

func New(l Locker) *Cond {
	mu := make(chan struct{}, 1)
	mu <- struct{}{}
	return &Cond{L: l, mu: mu}
}

func (c *Cond) Wait() {
	ch := make(chan struct{}, 1) // свой канал для этого горутина

	// Зарегистрироваться в списке
	<-c.mu
	c.waiters = append(c.waiters, ch)
	c.mu <- struct{}{}

	c.L.Unlock()
	<-ch // ждём именно своего сигнала
	c.L.Lock()
}

func (c *Cond) Signal() {
	<-c.mu
	if len(c.waiters) > 0 {
		ch := c.waiters[0]
		c.waiters = c.waiters[1:]
		c.mu <- struct{}{}
		ch <- struct{}{} // будим первого в очереди
	} else {
		c.mu <- struct{}{}
	}
}

func (c *Cond) Broadcast() {
	<-c.mu
	waiters := c.waiters
	c.waiters = nil
	c.mu <- struct{}{}
	for _, ch := range waiters {
		ch <- struct{}{} // будим всех
	}
}
