package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kr/pty"
)

func main() {

	c := exec.Command("/bin/bash", "-l")
	c.Env = append(os.Environ(), "TERM=xterm")
	tty, _ := pty.Start(c)
	sr := strings.NewReader("ls\n")

	ch := make(chan int, 0)

	go func() {
		for {
			read, _ := io.Copy(tty, sr)
			if read > 0 {
				ch <- 1
				fmt.Println("1111111111111111111111111")
				break
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ch:
				time.Sleep(1 * time.Second)
				buf := make([]byte, 1024)
				tty.Read(buf)
				fmt.Println("22222222222222222222222222")
				fmt.Println(string(buf))
			case <-time.After(2 * time.Second):
				continue
			}
		}

	}()

	for {

	}
}
