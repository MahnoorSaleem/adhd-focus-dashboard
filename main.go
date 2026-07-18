package main

import (
	"fmt"
	"io"
	"net/http"
)

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	response, err := io.ReadAll(r.Body)
	if err != nil {
		panic("something went wrong")
	}
	w.Write(response)
}

func main() {
	fmt.Println("Server Starting")
	http.HandleFunc("/dump", dumpHandler)
	http.ListenAndServe(":8080", nil)
}
