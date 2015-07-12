package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Job struct {
	Name  string
	Delay time.Duration
}

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
}

func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		fmt.Println("Starting worker:", i+1)
		worker := NewWorker(i+1, d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				fmt.Println("Adding job to workerJobQueue")
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}()
		}
	}
}

// NewWorker creates, and returns a new Worker object.
// It takes an id for identification and a reference to the
// worker pool from the dispatcher.
func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) start() {
	go func() {
		for {
			// Add my job queue into the worker pool.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				// Wait for dispatcher to add job to jobQueue
				fmt.Printf("worker%d: Received new job, delaying for %f seconds\n", w.id, job.Delay.Seconds())
				time.Sleep(job.Delay)
				fmt.Printf("worker%d: Hello, %s!\n", w.id, job.Name)
			case <-w.quitChan:
				// We have been asked to stop
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

func requestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
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

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Now, we retrieve the person's name from the request.
	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	// Now, we take the delay, and the person's name, and make a Job out of them.
	job := Job{Name: name, Delay: delay}

	// Push the work onto the queue.
	jobQueue <- job
	fmt.Println("Adding job to jobQueue")

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
}

func main() {
	// Parse the command-line flags.
	maxWorkers := flag.Int("max_workers", 5, "The number of workers to start")
	maxQueueSize := flag.Int("max_queue_size", 100, "The size of job queue")
	flag.Parse()

	// Create the job queue.
	jobQueue := make(chan Job, *maxQueueSize)

	// Start the dispatcher.
	dispatcher := NewDispatcher(jobQueue, *maxWorkers)
	dispatcher.run()

	// Start the HTTP handler.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue)
	})
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
