package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kr/pty"
)

func main() {

	c := exec.Command("/bin/bash", "-l")
	c.Env = append(os.Environ(), "TERM=xterm")
	tty, _ := pty.Start(c)
	str := strings.NewReader("top\n")
	io.Copy(tty, str)
	//ch := make(chan int, 0)

	go func() {
		for {

			buf := make([]byte, 1024*10)
			//buf := bytes.NewReader(nil)
			read, _ := tty.Read(buf)
			fmt.Println(read)
			fmt.Println(string(buf))
		}
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ch:
	// 			time.Sleep(1 * time.Second)
	// 			buf := make([]byte, 1024)
	// 			tty.Read(buf)
	// 			fmt.Println("22222222222222222222222222")
	// 			fmt.Println(string(buf))
	// 		case <-time.After(2 * time.Second):
	// 			continue
	// 		}
	// 	}

	// }()

	for {

	}
}
