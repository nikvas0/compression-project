// based on https://github.com/VKCOM/joy4/blob/master/examples/http_flv_and_rtmp_server/main.go

package rtmp

import (
	"log"
	"sync"

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
