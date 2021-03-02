package controller

import (
	"github.com/nikvas0/compression-project/video_service/rtmp"
)

// 1) Запускает rtmp cервер для трансляции
// 2) Запускает ffmpeg для конвертации в hls. В качестве входа использует rtmp сервер из п.1.
// 3) ffmpeg пишет в файлы, которые уже раздаются через nginx (ngx_http_hls_module)

// Запуск и окончание управляются через http: /start_stream/<path>, /stop_stream/<path>, /check_stream/<path>
// В п.2 может быть сразу несколько конвертаций для разных битрейтов/разрешений

type Controller struct {
	rtmpServer *rtmp.RtmpServer
}

func CreateController(rtmpServer *rtmp.RtmpServer) *Controller {
	controller := &Controller{rtmpServer}
	return controller
}

func (c *Controller) Run() {
	c.rtmpServer.Run()
}

func (c *Controller) StartStream(path string) {
	c.rtmpServer.AddStream("/" + path)
}

func (c *Controller) StopStream(path string) {
	c.rtmpServer.DeleteStream("/" + path)
}

func (c *Controller) CheckStream(path string) bool {
	return c.rtmpServer.StreamExists("/" + path)
}
