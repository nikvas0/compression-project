package main

import (
	"flag"
	"fmt"
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

func BasicVideo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	video := "http://localhost:8082/hls/" + r.URL.Query()["path"][0] + ".m3u8"

	fmt.Fprintf(w, `
	<html>
	<script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
	<video id="video" muted autoplay></video>
	<script>
	var video = document.getElementById('video');
	if(Hls.isSupported()) {
		var hls = new Hls();
		hls.loadSource('%s');
		hls.attachMedia(video);
		hls.on(Hls.Events.MANIFEST_PARSED,function() {
		video.play();
	});
	}
	else if (video.canPlayType('application/vnd.apple.mpegurl')) {
		video.src = '%s';
		video.addEventListener('loadedmetadata',function() {
		video.play();
		});
	}
	</script>
	</html>
	`, video, video)
}

func main() {
	log.Println("STARTING")

	port := flag.Int("port", 7191, "port")
	flag.Parse()

	http.HandleFunc("/video", BasicVideo)
	// http.HandleFunc("/create_stream", CreateStream)
	// http.HandleFunc("/get_stream", GetStream)
	// http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// })

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
