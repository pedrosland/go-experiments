# Job queues in Golang

A running example of the code from:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

I made a few adjustments to the code:

* Using non-exported private methods
* Removing the global variables
* Bringing the flags closer to their usage in main()

## Usage

Run the application:

```
PORT=5000 go run main.go -max_workers 5
```

Use flags to adjust the `max_workers` and `max_queue_size` when running the program.

Curl from another terminal:

```
for i in {1..15}; do curl localhost:5000/work -d name=$USER$i -d delay=$(expr $i % 5)s; done
```