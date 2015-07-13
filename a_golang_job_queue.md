# Golang Job Queue

A running example of the code from:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

I made a few adjustments to the code:

Use non-exported private methods
Remove global variables
Bring the flags closer to their usage in main()
Usage

## Run the application

    $ PORT=5000 go run example_dispatcher.go -max_workers 5

Use flags to adjust the max_workers and max_queue_size to override the default values.

Curl from another terminal window:

    $ for i in {1..15}; do curl localhost:5001/work -d name=$USER$i -d delay=$(expr $i % 9 + 1)s; done

## Remove Dispatcher

A simplified example that skips the dispatcher and creates workers directly. See `example_worker_only.go` for code.

## Performance

From what I can tell using Pprof the performance characteristics seem to remain the same between both examples.
