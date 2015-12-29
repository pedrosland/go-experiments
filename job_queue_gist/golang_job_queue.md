# Golang Job Queue

A running example of the code from:

* http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang
* http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

### Step 1

Small refactorings made to original code:

* Use non-exported private methods
* Remove global variables
* Bring the flags closer to their usage in main()

### Step 2

Simplify the worker queue by removing the `Dispatcher`.

* Creates workers directly and passes job queue to them

https://gist.github.com/harlow/dbcd639cf8d396a2ab73#file-worker_refactored-go

### Performance

The test run with Pprof show performance characteristics remain the same between both examples.

## Run the Application

Boot either the `worker_original.go` or the `worker_refactored.go` applications. Use flags to adjust the `max_workers` and `max_queue_size` to override the default values.

    $ go run worker_original.go -max_workers 5

cURL the application from another terminal window:

    $ for i in {1..15}; do curl localhost:8080/work -d name=job$i -d delay=$(expr $i % 9 + 1)s; done
