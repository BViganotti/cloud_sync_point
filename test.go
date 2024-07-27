package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Endpoint URL
	url := "http://127.0.0.1:3030/wait-for-second-party/F331860E-7745-49B4-8A0F-74249E954658" // generate the uuid dynamically

	// some data
	var jsonStr = []byte(`{"Hire_me":"Please"}`)

	// fire first post req
	sendPostRequest(url, jsonStr)

	time.Sleep(5 * time.Second)

	// fire second req
	sendPostRequest(url, jsonStr)
}

func sendPostRequest(url string, jsonStr []byte) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", body)
}
