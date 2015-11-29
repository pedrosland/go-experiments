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

// Job holds the attributes needed to perform unit of work.
type Job struct {
	Name  string
	Delay time.Duration
}

func newWorker(id int, jobQueue chan Job) worker {
	return worker{
		id:       id,
		jobQueue: jobQueue,
		quitChan: make(chan bool),
	}
}

type worker struct {
	id       int
	jobQueue chan Job
	quitChan chan bool
}

func (w worker) start() {
	go func() {
		for {
			select {
			case job := <-w.jobQueue:
				// Dispatcher has added a job to my jobQueue.
				fmt.Printf("worker%d: started %s, blocking for %f seconds\n", w.id, job.Name, job.Delay.Seconds())
				time.Sleep(job.Delay)
				fmt.Printf("worker%d: completed %s!\n", w.id, job.Name)
			case <-w.quitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

func (w worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

func requestHandler(jobQueue chan Job, w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the delay.
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate delay is in range 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Set name and validate value.
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// Create Job and push the work onto the jobQueue.
	job := Job{Name: name, Delay: delay}
	go func() {
		jobQueue <- job
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

	// Create the job queue.
	jobQueue := make(chan Job, *maxQueueSize)

	// Create new workers with access to jobQueue
	for i := 0; i < *maxWorkers; i++ {
		worker := newWorker(i, jobQueue)
		worker.start()
	}

	// Start the HTTP handler.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(jobQueue, w, r)
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
