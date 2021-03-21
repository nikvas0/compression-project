package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Describes Server.
type Server struct {
	streams []string
}

func CreateServer() *Server {
	server := &Server{}
	return server
}

func (s *Server) ListenAndServe(port string) {
	http.HandleFunc("/home", s.Home)
	http.HandleFunc("/start_stream", s.StartStream)
	http.HandleFunc("/stop_stream", s.StopStream)
	http.HandleFunc("/streams", s.Streams)
	http.HandleFunc("/video", s.Video)
	http.HandleFunc("/stream", s.Stream) // TODO

	log.Println("listening in port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	pageHtml, err := ioutil.ReadFile("../templates/home.html")
	fmt.Println(err)
	if err != nil {
		badRequest(w, "error reading home.html")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(pageHtml))
}

func (s *Server) StartStream(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["path"]) == 0 {
		badRequest(w, "path param not found")
		return
	}
	path := r.URL.Query()["path"][0]
	resp, err := http.Get("http://video_service/start_stream?path=" + path)
	if err != nil || resp.StatusCode != http.StatusOK {
		badRequest(w, "error requesting video_service start_stream")
		return
	}
	s.streams = append(s.streams, path)

	pageHtml, err := ioutil.ReadFile("../templates/message.html")
	fmt.Println(err)
	if err != nil {
		badRequest(w, "error reading message.html")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(pageHtml), "Stream was started successfully.")
}

func (s *Server) StopStream(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["path"]) == 0 {
		badRequest(w, "path param not found")
		return
	}
	path := r.URL.Query()["path"][0]
	resp, err := http.Get("http://video_service/stop_stream?path=" + path)
	if err != nil || resp.StatusCode != http.StatusOK {
		badRequest(w, "error requesting video_service stop_stream")
		return
	}
	var newStreams []string
	for _, stream := range s.streams {
		if stream != path {
			newStreams = append(newStreams, stream)
		}
	}
	s.streams = newStreams

	pageHtml, err := ioutil.ReadFile("../templates/message.html")
	fmt.Println(err)
	if err != nil {
		badRequest(w, "error reading message.html")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(pageHtml), "Stream was stopped successfully.")
}

func (s *Server) Streams(w http.ResponseWriter, r *http.Request) {
	pageHtml, err := ioutil.ReadFile("../templates/streams.html")
	fmt.Println(err)
	if err != nil {
		badRequest(w, "error reading streams.html")
		return
	}

	streamsTable := "<table>"
	for _, stream := range s.streams {
		streamsTable += "<tr>"
		streamsTable += fmt.Sprintf(`<td>%s</td>`, stream)
		streamsTable += fmt.Sprintf(`<td><a href="http://localhost:7191/video?path=%s&r=360p">360p (400k)</a></td>`, stream)
		streamsTable += fmt.Sprintf(`<td><a href="http://localhost:7191/video?path=%s&r=720p">720p (800k)</a></td>`, stream)
		streamsTable += fmt.Sprintf(`<td>http://localhost:8082/hls/%s_360p.m3u8</td>`, stream)
		streamsTable += "</tr>"
	}
	streamsTable += "</table>"

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(pageHtml), streamsTable)
}

func (s *Server) Video(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	get_video := func(resolution string) string {
		return "http://localhost:8082/hls/" + r.URL.Query()["path"][0] + "_" + resolution + ".m3u8"
	}
	video := get_video(r.URL.Query()["r"][0])

	get_page := func(resolution string) string {
		return "http://localhost:7191/video?path=" + r.URL.Query()["path"][0] + "&r=" + resolution
	}

	// unmute https://stackoverflow.com/a/39042127
	pageHtml, err := ioutil.ReadFile("../templates/video.html")
	fmt.Println(err)
	if err != nil {
		badRequest(w, "error reading video.html")
		return
	}
	fmt.Fprintf(w, string(pageHtml), video, video, get_page("360p"), get_page("720p"), r.URL.Query()["path"][0])
}

func (s *Server) Stream(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func badRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	writeErrorBody(w, err)
}

func writeErrorBody(w http.ResponseWriter, err interface{}) {
	var sErr string
	if err != nil {
		sErr = fmt.Sprint("error: ", err)
	} else {
		sErr = fmt.Sprint("Error while handling request.")
	}
	log.Println(sErr)
	w.Write([]byte(sErr))
}
