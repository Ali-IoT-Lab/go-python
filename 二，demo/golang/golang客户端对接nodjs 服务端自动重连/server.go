package main

import (
	"fmt"
	"time"

	"github.com/LightKool/socketio-client-go"
)

func main() {
	s, err := socketio.Socket("ws://127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	s.Connect()
	s.On("message", func(args ...interface{}) {
		fmt.Println("servver message!")
		fmt.Println(args[0])
	})
	for {
		s.Emit("messgae", "hello server!")
		time.Sleep(2 * time.Second)
	}
}
