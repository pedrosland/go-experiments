package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type job struct {
	name     string
	duration time.Duration
}

func newWorker(id int) *worker {
	return &worker{
		id:       id,
		quitChan: make(chan bool),
	}
}

type worker struct {
	id       int
	quitChan chan bool
}

func (w worker) process(jobCh chan job) {
	for {
		select {
		case j := <-jobCh:
			fmt.Printf("worker %d: started %s, duration: %f seconds\n", w.id, j.name, j.duration.Seconds())
			time.Sleep(j.duration)
		case <-w.quitChan:
			fmt.Printf("worker %d: stopped\n", w.id)
			return
		}
	}
}

func (w worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

func requestHandler(jobCh chan job, w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the durations.
	duration, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate delay is in range 1 to 10 seconds.
	if duration.Seconds() < 1 || duration.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Set name and validate value.
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// Create Job and push the work onto the jobCh.
	job := job{name, duration}
	go func() {
		fmt.Printf("added: %s %s\n", job.name, job.duration)
		jobCh <- job
	}()

	// Render success.
	w.WriteHeader(http.StatusCreated)
	return
}

func main() {
	var (
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		port         = flag.String("port", "8080", "The server port")
	)
	flag.Parse()

	// create job channel
	jobCh := make(chan job, *maxQueueSize)

	// create workers
	for i := 0; i < *maxWorkers; i++ {
		worker := newWorker(i)
		go worker.process(jobCh)
	}

	// handler for adding jobs
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(jobCh, w, r)
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
