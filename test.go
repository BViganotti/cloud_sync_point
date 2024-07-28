package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func make_post_req(url string, data []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("error making request:", err)
		return
	}
	defer resp.Body.Close()
	// read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:", err)
		return
	}
	// print the response
	fmt.Println("response:", string(body))
}

func main() {
	url := "http://127.0.0.1:3030/wait-for-second-party/F331860E-7745-49B4-8A0F-74249E954658"
	data := []byte(`{"Hire_me":"Please"}`)

	// we fire the second POST request in a different goroutine to not block the script
	go func() {
		fmt.Println("firing second POST request in 5s")
		time.Sleep(5 * time.Second)
		make_post_req(url, data)
	}()

	// firing the first POST request to the URI
	fmt.Println("sending first POST request")
	make_post_req(url, data)

	// wait to ensure the second request completes before the script exits
	time.Sleep(6 * time.Second)
}
