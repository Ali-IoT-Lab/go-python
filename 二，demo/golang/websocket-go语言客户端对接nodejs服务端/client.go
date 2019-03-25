package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	t := time.Now()
	fmt.Printf("当前的时间是: %d-%d-%d %d:%d:%d\n", t.Year(),
		t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl("127.0.0.1", 3000, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		t1 := time.Now()
		fmt.Printf("当前的时间是: %d-%d-%d %d:%d:%d\n", t1.Year(),
			t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second())
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {

		h.Emit("messgae", "asdasdasd")
		log.Println("Connected")
	})
	err = c.On("data", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	for {

	}
}
