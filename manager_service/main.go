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

	get_video := func(resolution string) string {
		return "http://localhost:8082/hls/" + r.URL.Query()["path"][0] + "_" + resolution + ".m3u8"
	}
	video := get_video(r.URL.Query()["r"][0])

	get_page := func(resolution string) string {
		return "http://localhost:7191/video?path=" + r.URL.Query()["path"][0] + "&r=" + resolution
	}

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

	<table>
	<tr>
		<td><a href="%s">360p (400k)</a>  </td>
		<td><a href="%s">720p (800k)</a>  </td>
	</tr>
	</table>
	</html>
	`, video, video, get_page("360p"), get_page("720p"))
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
