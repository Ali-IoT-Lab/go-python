package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/Ali-IoT-Lab/socketio-client-go"
	"github.com/kr/pty"
)

type wsPty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

func (wp *wsPty) Start() {
	var err error
	args := flag.Args()
	wp.Cmd = exec.Command(cmdFlag, args...)
	wp.Cmd.Env = append(os.Environ(), "TERM=xterm")
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

var cmdFlag string

func init() {
	flag.StringVar(&cmdFlag, "cmd", "/bin/bash", "command to execute on slave side of the pty")
	// TODO: make sure paths exist and have correct permissions
}
func main() {

	wp := wsPty{}
	// TODO: check for errors, return 500 on fail
	wp.Start()

	var Header http.Header = map[string][]string{
		"moja":     {"ccccc, asdasdasdasd"},
		"terminal": {"en-esadasdasdwrw"},
		"success":  {"dasdadas", "wdsadaderew"},
	}

	s, err := socketio.Socket("ws://127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
	s.Connect(Header)
	s.Emit("messgae", "hello server!")

	s.On("message", func(args ...interface{}) {
		fmt.Println("-------------------11111-------------------------")
		//fmt.Printf("%T\n", args[0])
		res, _ := GetBytes(args[0])
		fmt.Println(string(res))
		////_, ok := args.(string)
		wp.Pty.Write(res)
	})

	go func() {
		resBuf := make([]byte, 1024)
		// TODO: more graceful exit on socket close / process exit
		for {
			fmt.Println("-----------------------2222222---------------------------")
			fmt.Println(string(resBuf))

			n, err := wp.Pty.Read(resBuf)
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				return
			}
			fmt.Println(n)
			// 	out := make([]byte, base64.StdEncoding.EncodedLen(n))
			// 	base64.StdEncoding.Encode(out, resBuf[0:n])

			// 	fmt.Println("-----------------------3333333---------------------------")
			// 	fmt.Println(string(resBuf[0:n]))

			// 	s.Emit("result", string(resBuf[0:n]))
		}

	}()

	time.Sleep(2 * time.Second)

	//wp.Stop()
	for {
		// s.Emit("messgae", "hello server!")
		// time.Sleep(2 * time.Second)
	}
}
