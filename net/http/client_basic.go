package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	client := &http.Client{}
	url := fmt.Sprintf("http://%s", addr)
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

		// This body need to be read to the end and close, to let client reuse socket underline.
		respData, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalln("read response body failed:", err)
		}
		fmt.Println("return:", resp.StatusCode, resp.Status)
		fmt.Println("body:", string(respData))
	}
}
