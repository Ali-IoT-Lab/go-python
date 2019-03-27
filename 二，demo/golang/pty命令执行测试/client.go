




package main

import (
	"github.com/kr/pty"
	"io"
	"os"
	"os/exec"
)

func main() {
	c :=  exec.Command("/bin/bash", "-l")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	go func() {
		f.Write([]byte("l"))
		
		f.Write([]byte("s"))
		f.Write([]byte("\r"))
                f.Write([]byte{4}) // EOT
	}()
	io.Copy(os.Stdout, f)
}
