package main

import (
	"fmt"
	"time"
)

func main() {
	<-time.After(time.Second * time.Duration(5)) //定时2s，阻塞2s, 2s后产生一个事件，往channel写内容
	fmt.Println("时间到123123")
}

func main02() {
	time.Sleep(2 * time.Second)
	fmt.Println("时间到")
}

func main01() {
	//延时2s后打印一句话
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	fmt.Println("时间到")
}
