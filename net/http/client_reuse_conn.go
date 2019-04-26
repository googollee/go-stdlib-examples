package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const addr = "localhost:8888"

func echo(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(w, r.Body); err != nil {
		log.Println("echo request failed:", err)
	}
}

func server() {
	if err := http.ListenAndServe(addr, http.HandlerFunc(echo)); err != nil {
		log.Fatalln("can't listen", addr, ":", err)
	}
}

func main() {
	go server()

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			fmt.Println("dial a new connection")
			// Always with a timeout to make sure no hunging.
			return net.DialTimeout(network, addr, time.Second)
		},
	}
	client := &http.Client{
		Transport: transport,
	}
	url := fmt.Sprintf("http://%s", addr)

	fmt.Println("when read all from response.Body")
	for i := 0; i < 10; i++ {
		// Prepare request
		data := strings.NewReader("data")
		req, err := http.NewRequest("POST", url, data)
		if err != nil {
			log.Fatalln("create request error:", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("post to", url, "failed:", err)
		}

		// This body need to be read to the end, to let client reuse socket underline.
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("read response body failed:", err)
		}
		fmt.Println("return:", resp.StatusCode, resp.Status)
		fmt.Println("body:", string(respData))
	}

	fmt.Println("when not read all from response.Body")
	for i := 0; i < 10; i++ {
		// Prepare request
		data := strings.NewReader("data")
		req, err := http.NewRequest("POST", url, data)
		if err != nil {
			log.Fatalln("create request error:", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln("post to", url, "failed:", err)
		}

		fmt.Println("return:", resp.StatusCode, resp.Status)
	}
}
