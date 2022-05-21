package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	const (
		maxWorkers    = 4
		maxQueuedJobs = 20
		port          = ":8081"
	)

	jobQueue := make(chan Job, maxQueuedJobs)
	dispatcher := NewDispatcher(jobQueue, maxWorkers)

	dispatcher.Run()

	http.HandleFunc("/fib", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue)
	})
	log.Fatal(http.ListenAndServe(port, nil))
}

func requestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Invalid delay", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(r.FormValue("number"))
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	job := Job{Name: name, Delay: delay, Number: number}
	jobQueue <- job
	w.WriteHeader(http.StatusCreated)
}
