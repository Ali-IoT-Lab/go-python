package main

import "fmt"

type Person struct {
	name string
	sex  byte
	age  int
}

//指针作为接收者，引用语义
func (p *Person) SetInfoPointer() {
	(*p).name = "yoyo"
	(*p).sex = 'f'
	(*p).age = 22
	fmt.Println("SetInfoPointer: 	(*p) = ", (*p))
}

//值作为接收者，值语义
func (p Person) SetInfoValue() {
	p.name = "xxx"
	p.sex = 'm'
	p.age = 33

	fmt.Println("SetInfoValue: p = ", p)
}

func main() {
	//p 为指针类型
	var p *Person = &Person{"mike", 'm', 18}
	p.SetInfoPointer() //func (p) SetInfoPointer()

	p.SetInfoValue()    //func (*p) SetInfoValue()
	(*p).SetInfoValue() //func (*p) SetInfoValue()
}
