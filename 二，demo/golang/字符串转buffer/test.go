package main

import (
	"bytes"
	"fmt"
)

func main() {

	buf := bytes.NewBufferString("hello")
	s := buf.Bytes()
	fmt.Println(s)         //buf转整形
	fmt.Println(string(s)) //buf转整形

}
