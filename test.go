package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func regular_test(url string, data []byte) {
	// this test tests the regular usage of the webservice, a first request is send ->
	// the webservice blocks until the second request arrives.
	id := uuid.New()
	urlWithUUID := fmt.Sprintf("%s/%s", url, id.String())

	// we fire the second POST request in a different goroutine to not block the script
	go func() {
		fmt.Println("firing second POST request in 5s")
		time.Sleep(5 * time.Second)
		make_post_req(urlWithUUID, data)
	}()

	// firing the first POST request to the URI
	fmt.Println("sending first POST request")
	make_post_req(urlWithUUID, data)
}

func timeout_test(url string, data []byte) {
	// this test will on purpose send the second request 15seconds after the first one to
	// trigger the TIMEOUT response from the webservice.
	id := uuid.New()
	urlWithUUID := fmt.Sprintf("%s/%s", url, id.String())

	// we fire the second POST request in a different goroutine to not block the script
	go func() {
		fmt.Println("firing second POST request in 15s")
		time.Sleep(15 * time.Second)
		make_post_req(urlWithUUID, data)
	}()

	// firing the first POST request to the URI
	fmt.Println("sending first POST request")
	make_post_req(urlWithUUID, data)
}

func different_uuids_test(url string, data []byte) {
	// this test send two requests with different uuids to see how the
	// webservice reacts, they should both timeout
	fid := uuid.New()
	firstUrlWithUUID := fmt.Sprintf("%s/%s", url, fid.String())

	// we fire the second POST request in a different goroutine to not block the script
	go func() {
		sid := uuid.New()
		secondUrlWithUUID := fmt.Sprintf("%s/%s", url, sid.String())
		fmt.Println("firing second POST request in 2s")
		time.Sleep(2 * time.Second)
		make_post_req(secondUrlWithUUID, data)
	}()

	// firing the first POST request to the URI
	fmt.Println("sending first POST request")
	make_post_req(firstUrlWithUUID, data)
}

func main() {
	url := "http://127.0.0.1:3030/wait-for-second-party"
	data := []byte(`{"hire me":"please"}`)

	regular_test(url, data)
	timeout_test(url, data)
	different_uuids_test(url, data)

	// wait to ensure the second request completes before the script exits
	time.Sleep(6 * time.Second)
}
