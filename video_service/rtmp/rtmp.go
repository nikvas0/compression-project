// based on https://github.com/VKCOM/joy4/blob/master/examples/http_flv_and_rtmp_server/main.go

package rtmp

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/VKCOM/joy4/av/avutil"
	"github.com/VKCOM/joy4/av/pubsub"
	"github.com/VKCOM/joy4/format/rtmp"
)

type Channel struct {
	q         *pubsub.Queue
	hasHeader bool
}

type RtmpServer struct {
	queues map[string]*Channel
	server *rtmp.Server
	mux    *sync.RWMutex
}

func (s *RtmpServer) AddStream(path string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.queues[path]
	if ok {
		return
	}
	s.queues[path] = &Channel{pubsub.NewQueue(), false}
	log.Printf("Created stream: %s\n", path)
}

func (s *RtmpServer) StreamExists(path string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, ok := s.queues[path]
	return ok
}

func (s *RtmpServer) DeleteStream(path string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.queues[path]
	if !ok {
		return
	}
	delete(s.queues, path)
	log.Printf("Deleted stream: %s\n", path)
}

func (s *RtmpServer) Run() error {
	return s.server.ListenAndServe()
}

func CreateRtmpServer(addr string) *RtmpServer {
	if addr == "" {
		addr = ":1935"
	}

	s := &RtmpServer{map[string]*Channel{}, &rtmp.Server{}, &sync.RWMutex{}}
	s.server.Addr = addr

	s.server.HandlePlay = func(conn *rtmp.Conn) {
		s.mux.Lock()
		c, ok := s.queues[conn.URL.Path]
		s.mux.Unlock()
		if !ok || c == nil {
			log.Printf("Not found stream %s\n", conn.URL)
			return
		}
		if !c.hasHeader {
			log.Printf("Not found header in stream %s\n", conn.URL)
			return
		}
		err := avutil.CopyFile(conn, c.q.Latest())
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			return
		}
	}

	s.server.HandlePublish = func(conn *rtmp.Conn) {
		streams, err := conn.Streams()
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			return
		}

		s.mux.Lock()
		c, ok := s.queues[conn.URL.Path]
		if !ok || c == nil {
			log.Printf("Not found stream %v\n", conn.URL.Path)
			s.mux.Unlock()
			return
		}
		if !c.hasHeader {
			c.q.WriteHeader(streams)
			c.hasHeader = true
		}
		s.mux.Unlock()

		// TODO: переделать
		go func() {
			time.Sleep(1 * time.Second)
			cmd := exec.Command("ffmpeg", "-v", "verbose", "-i", "rtmp://localhost:1935"+conn.URL.Path,
				"-vf", "scale=w=640:h=360:force_original_aspect_ratio=decrease", "-c:a", "aac", "-ac", "1", "-c:v", "libx264", "-profile:v", "main", "-crf", "18",
				"-b:v", "400k", "-maxrate", "400k", "-bufsize", "80k", "-b:a", "80k",
				"-max_muxing_queue_size", "1024", "-flags", "-global_header", "-hls_time", "1", "-hls_list_size", "1", "-hls_wrap", "3", "-start_number", "1", "hls"+conn.URL.Path+"_360p.m3u8",
				"-vf", "scale=w=1280:h=720:force_original_aspect_ratio=decrease", "-c:a", "aac", "-ac", "1", "-c:v", "libx264", "-profile:v", "main", "-crf", "18",
				"-b:v", "800k", "-maxrate", "800k", "-bufsize", "80k", "-b:a", "80k",
				"-max_muxing_queue_size", "1024", "-flags", "-global_header", "-hls_time", "1", "-hls_list_size", "1", "-hls_wrap", "3", "-start_number", "1", "hls"+conn.URL.Path+"_720p.m3u8")

			//err := cmd.Start()
			var out bytes.Buffer
			var er bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &er
			err := cmd.Run()
			if err != nil {
				fmt.Printf("out: %q\n", out.String())
				fmt.Printf("err: %q\n", er.String())
				log.Fatal(err)
			}
			fmt.Printf("out: %q\n", out.String())
			fmt.Printf("err: %q\n", er.String())
		}()

		err = avutil.CopyPackets(c.q, conn)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			return
		}

		s.mux.Lock()
		defer s.mux.Unlock()
		delete(s.queues, conn.URL.Path)
		c.q.Close()
	}

	return s
}
