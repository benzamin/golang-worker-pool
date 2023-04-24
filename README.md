# Golang Worker Pool Sample Implementation

This is a sample implementation of a worker pool modality in Golang, which can be used to execute a set of jobs concurrently. The worker pool is designed to be highly generic and can be used by just tweaking the worker count and job queue limit. **Disclaimer:** I've *assembled* this piece of code as a part of my learning on how high concurrency system works in golang.

## Overview

The worker pool modality is a common design pattern used to execute a set of jobs concurrently by utilizing a pool of workers. The basic idea is to have a set of workers waiting for jobs to be added to a job queue, and once a job is added to the queue, a free worker will pick it up and execute it, and will return the result using a result channel. By using this approach, we can execute multiple jobs concurrently using the famous `goroutine`, and also limit the number of workers and the number of jobs in the queue, which can be useful in situations where we need to control the load on the system, a perfectly tweaked worker-pool implementation can handle millions of reuqests per second.

## Implementation 

The worker pool implementation in this repository consists of the following files:

    worker_pool.go: This file contains the main implementation of the worker pool. It defines a WorkerPool struct, which holds the job queue and the worker channels, and along with dispatcher.

    heavy_task.go: This shows how to implement a generic task which can be passed to worker pool. Currently it simulates a random sleep, but you can add anything here (please check the 'post' branch for an example).

    main.go: This file contains a simple example of how to use the worker pool. It creates a pool with 5 workers and a job queue limit of 100. The params can be tweaked while running.
    
## Usage

Clone this repo, then simply run this project:

    go run .

Or tweak the params while you run:

    go run . -max_workers 20 -max_queue_size 100 -port 8080

To use the worker pool in your own code, you can simply copy the worker_pool.go file to your project. Then you can create a new jobqueue and dispatcher object with the desired number of job queue limit and workers, and run the dispatcher. Adding jobs to the jobqueue will start concurrent processing.

The Job, JobResult and Task are simple struct which can be used generically.

    type Job struct {
        Task          Task
        Payload       Payload
        ReturnChannel chan JobResult
    }

    type JobResult struct {
        Value interface{}
        Error error
    }

The task must implementt a `Run` method that will be called from inside a worker routine.
    
    type HeavyTask struct{}

    func NewHeavyTask() HeavyTask {
        return HeavyTask{}
    }

    func (e *HeavyTask) Run(payload Payload) JobResult {
        //execute your task here
    }


### Here's an example of main function, which expose a GET endpoint, upon receiving a request it starts processing the job using the worker pool:

    func main() {
    ......
    ......
        jobQueue := make(chan Job, 10)

        // Start the dispatcher.
        dispatcher := NewDispatcher(jobQueue, 5)
        dispatcher.run()

        // Start the HTTP handler.
        http.HandleFunc("/heavyapi", func(w http.ResponseWriter, r *http.Request) {
            heavyApiHandlerGET(w, r, jobQueue)
        })
    }

This example creates a worker pool with 5 workers and a job queue limit of 10, and then start the dispatcher, waiting for job. Later a http endpoint called `/heavyapi` is exposed, which handler actually adds a new task in the jobque. It waits for a the result using a return channel, and process it once get the result. These jobs are executed concurrently by the workers in the pool.

    task := NewHeavyTask()
        parameters := map[string]string{"sleep": sleep, "random": random}
        job := NewJob(&task, parameters, NewJobResultChannel())
        jobQueue <- job

        resp := <-job.ReturnChannel
        if resp.Error != nil {
            sendFailResponse(w, http.StatusInternalServerError, fmt.Sprintf("Something went wrong. %s", resp.Error))
            return
        }
        sendSuccessResponse(w, resp.Value)

## Benchmarking
Some sample load test was done for this code using [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html). Here are some tests I've done, where -c=concurrenct requests to send, -n= number of total requests to send, -v=verbose level 2, -p=filename for post requeest body.

### General test 
    ab -c 50 -n 1000 -v 2 "http://localhost:8080/heavyapi?sleep=35" > log.txt

### Test with random sleep
    ab -c 50 -n 1000 -v 2 "http://localhost:8080/heavyapi?random=true" > log.txt

### To test POST endpoint (branch 'post'), create a file called "post.txt" and put the following line inside the file:
    {"product":"toyota"}

### Then run the POST test [depends on external dummy API]
    ab -c 10 -n 100 -v 2 -p post.txt "http://localhost:8080/bigapi" > log.txt

## Conclusion

The worker pool modality is a powerful tool for executing a set of jobs concurrently in Golang. By using a worker pool, we can control the load on the system and also improve the performance of our code. The implementation provided in this repository is highly generic and can be used in a variety of situations by just tweaking the worker count and job queue limit.
