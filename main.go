package main

import (
	"encoding/json"
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	var (
		maxWorkers   = flag.Int("max_workers", 5, "The number of workers to start")
		maxQueueSize = flag.Int("max_queue_size", 100, "The size of job queue")
		port         = flag.String("port", "8080", "The server port")
	)

	flag.Parse()

	fmt.Printf("Starting Worker-Pool using: \n> max_workers: %d \n> max_queue_size: %d \n> port: %s \n\n", *maxWorkers, *maxQueueSize, *port)

	// Create the job queue.
	jobQueue := make(chan Job, *maxQueueSize)

	// Start the dispatcher.
	dispatcher := NewDispatcher(jobQueue, *maxWorkers)
	dispatcher.run()
	logStats()

	// Start the HTTP handler.
	http.HandleFunc("/heavyapi", func(w http.ResponseWriter, r *http.Request) {
		heavyApiHandlerGET(w, r, jobQueue)
	})

	server := &http.Server{
		Addr:           ":" + *port,
		ReadTimeout:    time.Duration(5000) * time.Millisecond,
		WriteTimeout:   time.Duration(10000) * time.Millisecond,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())

}

func heavyApiHandlerGET(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
	// Make sure we can only be called with an HTTP GET request.
	if r.Method != "GET" {
		sendFailResponse(w, http.StatusMethodNotAllowed, "You must use GET method")
		return
	}

	sleep := r.URL.Query().Get("sleep")
	random := r.URL.Query().Get("random")

	if sleep == "" && random == "" {
		sendFailResponse(w, http.StatusBadRequest, "You must specify a search term, ex: /search?q=internet")
		return
	}

	// Create Job and push the work onto the jobQueue, wait for response
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
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func sendSuccessResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{true, "Success", response})
}

func sendFailResponse(w http.ResponseWriter, httpStatus int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(&Response{false, errorMessage, nil})
}
