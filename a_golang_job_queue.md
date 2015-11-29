# Golang Job Queue

A running example of the code from:

http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

### Step 1

I made a few adjustments to the code:

* Use non-exported private methods
* Remove global variables
* Bring the flags closer to their usage in main()
* Usage

### Step 2

Removing the `Dispatcher`. A simplified example that removes the dispatcher and creates workers directly.

https://gist.github.com/harlow/dbcd639cf8d396a2ab73#file-example_worker_only-go

## Run the Application

    $ PORT=5000 go run example_dispatcher.go -max_workers 5

Use flags to adjust the `max_workers` and `max_queue_size` to override the default values.

Curl from another terminal window:

    $ for i in {1..15}; do curl localhost:5001/work -d name=$USER$i -d delay=$(expr $i % 9 + 1)s; done

## Performance

From what I can tell using Pprof the performance characteristics seem to remain the same between both examples.
