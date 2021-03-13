package main

import (
	"flag"
	"fmt"
	"github.com/nikvas0/compression-project/chat_service/client"
	"log"
	"net/http"
	"strconv"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	room := r.URL.Query()["path"][0]
	user := r.URL.Query()["user"][0]
	message := r.URL.Query()["message"][0]
	fmt.Println(room, user, message)
	client.SendMessage(user, room, message)
}

func main() {
	log.Println("Starting chat handler")

	port := flag.Int("port", 7192, "port")
	flag.Parse()

	http.HandleFunc("/sendmessage", SendMessage)
	// http.HandleFunc("/create_stream", CreateStream)
	// http.HandleFunc("/get_stream", GetStream)
	// http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
