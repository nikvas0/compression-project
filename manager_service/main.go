package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func CreateStream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetStream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("STARTING")

	port := flag.Int("port", 7191, "port")
	flag.Parse()

	// http.HandleFunc("/create_stream", CreateStream)
	// http.HandleFunc("/get_stream", GetStream)
	// http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
