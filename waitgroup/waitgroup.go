//go:build !solution

package waitgroup

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
type WaitGroup struct {
	waitCh chan struct{}
	cntrCh chan int
}

// New creates WaitGroup.
func New() *WaitGroup {
	waitCh := make(chan struct{}, 1)
	cntrCh := make(chan int, 1)
	waitCh <- struct{}{}
	cntrCh <- 0
	return &WaitGroup{waitCh: waitCh, cntrCh: cntrCh}
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	n := <-wg.cntrCh
	if n == 0 && delta > 0 {
		<-wg.waitCh
	}
	n += delta
	if n < 0 {
		panic("negative WaitGroup counter")
	}
	wg.cntrCh <- n
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	n := <-wg.cntrCh
	n--
	if n < 0 {
		panic("negative WaitGroup counter")
	} else if n == 0 {
		wg.waitCh <- struct{}{}
	}
	wg.cntrCh <- n
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	if _, ok := <-wg.waitCh; ok {
		defer func() {
			wg.waitCh <- struct{}{}
		}()
	}
}
