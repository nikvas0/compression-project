package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/nikvas0/compression-project/chat_service/protobufs"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

var clients map[string]chat.BroadcastClient
var waits map[string]*sync.WaitGroup
var usersToId map[string]string
var idToUsers map[string]string

func init() {
	clients = make(map[string]chat.BroadcastClient)
	waits = make(map[string]*sync.WaitGroup)
	usersToId = make(map[string]string)
	idToUsers = make(map[string]string)
}

func connect(user *chat.User, room string) error {
	var streamerror error
	stream, err := clients[user.Id].CreateStream(context.Background(), &chat.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	waits[user.Id] = &sync.WaitGroup{}
	waits[user.Id].Add(1)
	go func(str chat.Broadcast_CreateStreamClient) {
		defer waits[user.Id].Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}
			if msg.Room == room {
				fmt.Printf("%v : %s\n", msg.User.DisplayName, msg.Message)
			}
		}
	}(stream)
	return streamerror
}

func ConnectUser(name string, room string) {
	timestamp := time.Now()
	id := sha256.Sum256([]byte(timestamp.String() + name))
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}
	user := &chat.User{
		Id:          hex.EncodeToString(id[:]),
		DisplayName: name,
	}
	usersToId[user.DisplayName] = user.Id
	idToUsers[user.Id] = user.DisplayName
	clients[user.Id] = chat.NewBroadcastClient(conn)
	connect(user, room)
}

func SendMessage(user string, room string, message string) {
	if _, ok := clients[usersToId[user]]; !ok {
		ConnectUser(user, room)
	}
	timestamp := time.Now()
	msgId := sha256.Sum256([]byte(timestamp.String() + user))
	msg := &chat.Message{
		Id: hex.EncodeToString(msgId[:]),
		User: &chat.User{
			Id:          usersToId[user],
			DisplayName: user,
		},
		Room:      room,
		Message:   message,
		Timestamp: timestamp.String(),
	}
	_, err := clients[usersToId[user]].BroadcastMessage(context.Background(), msg)
	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}
}

func main() {
	name := flag.String("N", "Guest", "")
	room := flag.String("R", "Default", "")
	flag.Parse()
	ConnectUser(*name, *room)
}