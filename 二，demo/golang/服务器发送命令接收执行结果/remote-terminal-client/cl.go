package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"os"
	"os/exec"

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
}

func main() {
	wp := wsPty{}
	wp.Start()
	go func() {
		resBuf := make([]byte, 1024)

		for {
			_, err := wp.Pty.Read(resBuf)
			println(string(resBuf))
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				return
			}
		}
	}()

	go func() {
		wp.Pty.Write([]byte("t"))
		wp.Pty.Write([]byte("o"))
		wp.Pty.Write([]byte("p"))
		wp.Pty.Write([]byte("\n"))

	}()

	//time.Sleep(2 * time.Second)

	for {
	}
}



------------------

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

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

	// go func() {
	// 	wp.Pty.Write([]byte("t"))
	// 	wp.Pty.Write([]byte("o"))
	// 	wp.Pty.Write([]byte("p"))
	// 	wp.Pty.Write([]byte("\n"))

	// }()

	go func() {

		s.On("message", func(args ...interface{}) {
			fmt.Println("--------------------------------------------------")

			//fmt.Printf("%T\n", args[0])
			fmt.Println(args[0])

			// cmd := args[0]
			// payload, _ := GetBytes(cmd)

			// cmdd := strings.Replace(string(payload), " ", "", -1)
			// fmt.Println(len(string(payload)))
			// buf := make([]byte, base64.StdEncoding.DecodedLen(len(payload)))
			// _, err := base64.StdEncoding.Decode(buf, payload)
			// if err != nil {
			// 	log.Printf("base64 decoding of payload failed: %s\n", err)
			// }
			// wp.Pty.Write(buf)
			//res, _ := GetBytes(args[0])
			//res = res[0:1]
			// fmt.Print(string(res))
			// fmt.Print(string(res[1:2]))

			//_, ok := args.(string)
			// wp.Pty.Write([]byte("l"))
			// wp.Pty.Write([]byte("s"))
			//wp.Pty.Write([]byte("p"))
			// wp.Pty.Write([]byte("\n"))
		})
	}()

	go func() {
		buf := make([]byte, 1024)
		// TODO: more graceful exit on socket close / process exit
		for {
			// fmt.Println("--------------------------------------------------")
			// fmt.Println(string(buf))

			n, err := wp.Pty.Read(buf)
			fmt.Println(n)
			if err != nil {
				log.Printf("Failed to read from pty master: %s", err)
				break
			}
			// fmt.Println("-----------------------4444444---------------------------")
			// fmt.Println(n)
			out := make([]byte, base64.StdEncoding.EncodedLen(n))
			base64.StdEncoding.Encode(out, buf[0:n])

			// fmt.Println("-----------------------3333333---------------------------")
			// fmt.Println(string(buf[0:n]))

			s.Emit("result", string(buf[0:n]))
		}

	}()

	//time.Sleep(2 * time.Second)

	//wp.Stop()
	for {
		// s.Emit("messgae", "hello server!")
		// time.Sleep(2 * time.Second)
	}
}
