package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kr/pty"
)

func main() {
	c := exec.Command("/bin/bash", "-l")
	f, err := pty.Start(c)
	c.Env = os.Environ()
	c.Env = append(c.Env, "LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8")
	if err != nil {
		panic(err)
	}

	go func() {
		f.Write([]byte("l"))
		f.Write([]byte("s"))
		f.Write([]byte("\r"))
		//f.Write([]byte{4}) // EOT

	}()

	time.Sleep(2 * time.Second)
	//var file *os.File
	//io.Copy(file, f)
	buf := make([]byte, 1024)
	f.Read(buf)
	var strf = string(buf)
	var a = strings.Split(strf, "\n")
	fmt.Println(a[2])
}
