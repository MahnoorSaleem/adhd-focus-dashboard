package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var latestResult string
var mu sync.Mutex
var jobQueue = make(chan string, 10) // holds up to 10 pending jobs

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	response, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// go func() {
	// 	result := fakeAIProcess((string(response)))
	// 	mu.Lock()
	// 	latestResult = result
	// 	mu.Unlock()
	// 	if err != nil {
	// 		panic("something went wrong")
	// 	}
	// }()
	jobQueue <- string(response)
	w.Write([]byte("Got it processing"))
}

func fakeAIProcess(text string) string {
	time.Sleep((3 * time.Second))
	return "Task: " + text

}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	result := latestResult
	mu.Unlock()

	w.Write(([]byte(result)))
}

func worker(id int) {
	for text := range jobQueue {
		fmt.Println("Worker", id, "Started at", time.Now().Format("15:04:05"), "-", text)
		result := fakeAIProcess(text)
		mu.Lock()
		latestResult = result
		mu.Unlock()
		fmt.Println("Worker", id, "Finished at", time.Now().Format("15:04:05"), "-", result)
	}

}

func main() {
	fmt.Println("Server Starting")
	http.HandleFunc("/dump", dumpHandler)
	http.HandleFunc("/latest", latestHandler)
	for i := 1; i <= 3; i++ {
		go worker(i)
	}
	http.ListenAndServe(":8080", nil)

}
