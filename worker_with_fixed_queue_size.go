package main

import (
	_ "expvar"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct {
	name     string
	duration time.Duration
}

type worker struct {
	id int
}

func (w worker) process(j job) {
	fmt.Printf("worker%d: started %s, working for %f seconds\n", w.id, j.Name, j.Delay.Seconds())
	time.Sleep(j.duration)
	fmt.Printf("worker%d: completed %s!\n", w.id, j.Name)
}

func main() {
	wg := &sync.WaitGroup{}
	jobCh := make(chan job)

	// start workers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		w := &worker{i}

		go func(w worker) {
			for j := range jobCh {
				w.process(j)
			}
			wg.Done()
		}(w)
	}

	// add jobs to queue
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("job-%d", i)
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		fmt.Printf("adding: %s %s\n", name, duration)
		jobCh <- job{name, duration}
	}

	// close jobCh and wait for workers to be done
	close(jobCh)
	wg.Wait()
}
