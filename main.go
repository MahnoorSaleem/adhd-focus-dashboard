package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	response, err := io.ReadAll(r.Body)
	result := FakeAIProcess((string(response)))
	if err != nil {
		panic("something went wrong")
	}
	w.Write([]byte(result))
}

func FakeAIProcess(text string) string {
	time.Sleep((3 * time.Second))
	return "Task: " + text

}

func main() {
	fmt.Println("Server Starting")
	http.HandleFunc("/dump", dumpHandler)
	http.ListenAndServe(":8080", nil)
}
