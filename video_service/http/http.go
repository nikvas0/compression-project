package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nikvas0/compression-project/video_service/controller"
)

type Server struct {
	controller *controller.Controller
}

func CreateServer(controller *controller.Controller) *Server {
	server := &Server{controller}
	return server
}

func (s *Server) Run() {
	s.controller.Run()
}

func (s *Server) ListenAndServe(port string) {
	http.HandleFunc("/start_stream", s.StartStream)
	http.HandleFunc("/stop_stream", s.StopStream)
	http.HandleFunc("/check_stream", s.CheckStream)

	go s.Run()
	log.Println("listening in port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) StartStream(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["path"]) == 0 {
		badRequest(w, "path param not found")
		return
	}
	path := r.URL.Query()["path"][0]
	s.controller.StartStream(path)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) StopStream(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["path"]) == 0 {
		badRequest(w, "path param not found")
		return
	}
	path := r.URL.Query()["path"][0]
	s.controller.StopStream(path)
	w.WriteHeader(http.StatusOK)
}

// CheckStream returns http.StatusOK on existing stream and http.StatusNotFound otherwise.
func (s *Server) CheckStream(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()["path"]) == 0 {
		badRequest(w, "path param not found")
		return
	}
	path := r.URL.Query()["path"][0]
	if s.controller.CheckStream(path) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
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
