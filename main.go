package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var resultChan = make(chan string)

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	response, err := io.ReadAll(r.Body)
	go func() {
		result := FakeAIProcess((string(response)))
		resultChan <- result
		if err != nil {
			panic("something went wrong")
		}
	}()
	w.Write([]byte("Git it processing"))
}

func FakeAIProcess(text string) string {
	time.Sleep((3 * time.Second))
	return "Task: " + text

}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	result := <-resultChan
	w.Write(([]byte(result)))
}

func main() {
	fmt.Println("Server Starting")
	http.HandleFunc("/dump", dumpHandler)
	http.HandleFunc("/latest", latestHandler)
	http.ListenAndServe(":8080", nil)
}
