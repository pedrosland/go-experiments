# Job queues in Golang

A running example of the code from:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

I made a few adjustments to the code:

* Use non-exported private methods
* Remove global variables
* Bring the flags closer to their usage in main()

## Usage

Run the application:

```
PORT=5000 go run main.go -max_workers 5
```

Use flags to adjust the `max_workers` and `max_queue_size` to override the default values.

Curl from another terminal window:

```
for i in {1..15}; do curl localhost:5000/work -d name=$USER$i -d delay=$(expr $i % 5)s; done
```

## Worker Only Queue

A simplified example that skips the dispatcher and creates workers directly. See `main_worker_only.go` the performance characteristics seem to remain the same.