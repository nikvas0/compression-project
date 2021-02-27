package main

import (
	"log"
	"net/http"
)

func CreateStream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetStream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	log.Println("STARTING")

	http.HandleFunc("/create_stream", CreateStream)
	http.HandleFunc("/get_stream", GetStream)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
