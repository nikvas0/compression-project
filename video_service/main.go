package main

import (
	"github.com/nikvas0/compression-project/video_service/rtmp"
)

func main() {
	// cм controller, который будет вызываться из http

	// rtmpToHls := converter.CreateConverter(...)

	rtmpServer := rtmp.CreateRtmpServer(":1935")
	rtmpServer.AddStream("/test")
	rtmpServer.Run()

	// controller := controller.CreateController(rtmp_server, rtmp_to_hls)

	// server := http.CreateServer(controller)
	// server.ListenAndServe()
}
