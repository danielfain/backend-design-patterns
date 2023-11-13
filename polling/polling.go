package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"
)

var finishedRequests = []int{}

func main() {
	pollType := os.Args[1]

	if pollType == "short" {
		ShortPollCheckRequest()
	} else {
		LongPollCheckRequest()
	}

	SubmitRequest()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SubmitRequest() {
	http.HandleFunc("/submitRequest", func(w http.ResponseWriter, r *http.Request) {
		requestId := time.Now().UnixMilli()
		log.Printf("Received request: %d", requestId)
		go ProcessRequest(int(requestId))
		fmt.Fprintf(w, "Received request: %d\n", requestId)
	})
}

func ShortPollCheckRequest() {
	http.HandleFunc("/checkRequest", func(w http.ResponseWriter, r *http.Request) {
		requestId, err := strconv.Atoi(r.URL.Query().Get("requestId"))

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Checking request with ID: %d\n", requestId)

		if slices.Contains(finishedRequests, requestId) {
			fmt.Fprintln(w, "Request finished!")
		} else {
			fmt.Fprintln(w, "Request still processing...")
		}

	})
}

func LongPollCheckRequest() {
	http.HandleFunc("/checkRequest", func(w http.ResponseWriter, r *http.Request) {
		requestId, err := strconv.Atoi(r.URL.Query().Get("requestId"))

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Checking request with ID: %d\n", requestId)

		isFinished := slices.Contains(finishedRequests, requestId)

		for !isFinished {
			isFinished = slices.Contains(finishedRequests, requestId)
			time.Sleep(1 * time.Second)
		}

		fmt.Fprintln(w, "Request finished!")
	})

}

func ProcessRequest(requestId int) {
	time.Sleep(15 * time.Second)
	finishedRequests = append(finishedRequests, requestId)
}
