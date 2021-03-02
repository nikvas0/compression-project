package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/nikvas0/compression-project/video_service/controller"
	"github.com/nikvas0/compression-project/video_service/http"
	"github.com/nikvas0/compression-project/video_service/rtmp"
)

func main() {
	// cм controller, который будет вызываться из http
	log.Println("STARTING")

	port := flag.Int("port", 80, "port")
	flag.Parse()

	// rtmpToHls := converter.CreateConverter(...)

	rtmpServer := rtmp.CreateRtmpServer(":1935")
	// rtmpServer.AddStream("/test")
	// rtmpServer.Run()

	controller := controller.CreateController(rtmpServer) // TODO доделать rtmp_to_hls

	server := http.CreateServer(controller)
	server.ListenAndServe(strconv.Itoa(*port))
}
