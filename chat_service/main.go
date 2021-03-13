package main

import (
	"context"
	"github.com/nikvas0/compression-project/chat_service/protobufs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
	"os"
	"sync"
)

var grpcLog grpclog.LoggerV2

func init() {
	grpcLog = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type Connection struct {
	stream chat.Broadcast_CreateStreamServer
	id     string
	active bool
	err    chan error
}

type Server struct {
	Connection []*Connection
}

func (s *Server) CreateStream(pconn *chat.Connect, stream chat.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		err:    make(chan error),
	}
	s.Connection = append(s.Connection, conn)
	return <-conn.err
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *chat.Message) (*chat.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *chat.Message, conn *Connection) {
			defer wait.Done()
			if conn.active {
				err := conn.stream.Send(msg)
				grpcLog.Infof("Sending message %v to user %v", msg.Id, conn.id)
				if err != nil {
					grpcLog.Errorf("Error with stream %v. Error: %v", conn.stream, err)
					conn.active = false
					conn.err <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &chat.Close{}, nil
}

func main() {
	var connections []*Connection

	server := &Server{connections}
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}
	chat.RegisterBroadcastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
