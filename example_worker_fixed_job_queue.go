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
	fmt.Printf("worker%d: doing %s (%s)\n", w.id, j.name, j.duration)
	time.Sleep(j.duration)
}

func startWorkers(jobCh chan job, wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		w := &worker{i}
		wg.Add(1)
		go func(w *worker) {
			for j := range jobCh {
				w.process(j)
			}
			wg.Done()
		}(w)
	}
}

func addJobs(jobCh chan job) {
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("job-%d", i)
		duration := time.Duration(rand.Intn(1000)) * time.Millisecond
		job := job{name, duration}
		fmt.Printf("adding: %s %s\n", job.name, job.duration)
		jobCh <- job
	}
	close(jobCh)
}

func main() {
	wg := &sync.WaitGroup{}
	jobCh := make(chan job)

	// start workers and add jobs
	startWorkers(jobCh, wg)
	addJobs(jobCh)

	// wait for workers to complete
	wg.Wait()
}
