package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/nikvas0/compression-project/manager_service/http"
)

func main() {
	log.Println("STARTING")

	port := flag.Int("port", 7191, "port")
	flag.Parse()

	// log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
	server := http.CreateServer()
	server.ListenAndServe(strconv.Itoa(*port))
}
