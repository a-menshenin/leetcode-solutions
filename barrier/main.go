package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Barrier struct {
	maxCount     atomic.Int32
	currentCount atomic.Int32
	c            *sync.Cond
}

func NewBarrier(n int) *Barrier {
	maxCount := atomic.Int32{}
	maxCount.Add(int32(n))
	currentCount := atomic.Int32{}

	return &Barrier{maxCount: maxCount, currentCount: currentCount, c: sync.NewCond(&sync.Mutex{})}
}

func (b *Barrier) Wait(id int) {
	b.c.L.Lock()
	b.c.Wait()
	b.c.L.Unlock()
}

func worker(id int, barrier *Barrier) {
	fmt.Printf("Горутина %d: начало работы до барьера\n", id)
	time.Sleep(time.Duration(id) * 500 * time.Millisecond) // симулируем работу
	fmt.Printf("Горутина %d: достигла барьера\n", id)

	barrier.currentCount.Store(barrier.currentCount.Add(1))

	if barrier.currentCount.Load() >= barrier.maxCount.Load() {
		fmt.Printf("id: %d : b.currentCount.Load() >= b.maxCount.Load()\n", id)
		barrier.c.Broadcast()
	} else {
		// Ждем на барьере, пока все не соберутся
		barrier.Wait(id)
	}

	fmt.Printf("Горутина %d: продолжение работы после барьера\n", id)
}

func main() {
	const numWorkers = 5
	barrier := NewBarrier(numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, barrier)
		}(i)
	}

	wg.Wait()
	fmt.Println("Все горутины завершили работу.")
}
