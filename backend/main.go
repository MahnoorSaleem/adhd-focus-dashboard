package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var latestResult string
var mu sync.Mutex
var jobQueue = make(chan string, 10) // holds up to 10 pending jobs
var groqAPIKey string

type groqRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type choice struct {
	Message message `json:"message"`
}

type groqResponse struct {
	Choices []choice `json:"choices"`
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mu.Lock()
	result := latestResult
	mu.Unlock()

	w.Write(([]byte(result)))
}

func worker(id int) {
	for text := range jobQueue {
		fmt.Println("Worker", id, "Started at", time.Now().Format("15:04:05"), "-", text)
		result := callAI(text)
		mu.Lock()
		latestResult = result
		mu.Unlock()
		fmt.Println("Worker", id, "Finished at", time.Now().Format("15:04:05"), "-", result)
	}

}

func callAI(text string) string {

	reqBody := groqRequest{
		Model:    "llama-3.1-8b-instant",
		Messages: []message{{Role: "user", Content: "Here is a messy brain dump. Give me only ONE clear, actionable task from it, nothing else: " + text}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "Error building request"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "Error creating request"
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+groqAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("REAL ERROR:", err)
		return "Error sending request"
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading response"
	}

	var groqResp groqResponse
	err = json.Unmarshal(body, &groqResp)
	if err != nil {
		return "Error parsing response"
	}

	return groqResp.Choices[0].Message.Content

}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

func main() {
	loadEnvFile(".env")
	groqAPIKey = os.Getenv("GROQ_KEY")
	if groqAPIKey == "" {
		fmt.Println("Warning: GROQ_KEY not set!")
	}
	fmt.Println("Server Starting")

	http.HandleFunc("/dump", dumpHandler)
	http.HandleFunc("/latest", latestHandler)
	for i := 1; i <= 3; i++ {
		go worker(i)
	}
	http.ListenAndServe(":8080", nil)

}
