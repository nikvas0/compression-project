package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nikvas0/compression-project/chat_service/client"
	"log"
	"net/http"
	"strconv"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "POST" {
		room := r.URL.Query()["path"][0]
		user := r.URL.Query()["user"][0]
		message := r.URL.Query()["message"][0]
		fmt.Println(room, user, message)
		client.SendMessage(user, room, message)
	}
	return
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		room := r.URL.Query()["path"][0]
		_ = json.NewEncoder(w).Encode(client.Messages[room])
		fmt.Println(client.Messages[room])
		client.Messages[room] = nil
	}
	return
}

func main() {
	log.Println("STARTING")

	port := flag.Int("port", 7192, "port")
	flag.Parse()

	http.HandleFunc("/sendmessage", SendMessage)
	http.HandleFunc("/getmessages", GetMessages)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
