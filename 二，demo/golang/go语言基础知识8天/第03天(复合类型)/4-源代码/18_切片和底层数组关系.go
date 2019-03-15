package main //必须有个main包

import "fmt"

func main() {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	//新切片
	s1 := a[2:5] //从a[2]开始，取3个元素
	fmt.Printf("s1=%p,cap=%d,len = %d %T \n", &s1, cap(s1), len(s1), s1)
	fmt.Printf("a=%p,cap=%d,len = %d %T\n", &a, cap(a), len(a), a)
	s1[1] = 666

	for i := 0; i < 1000000; i++ {
		a = append(a, 888)
		s1 = append(s1, 888)
	}

	fmt.Printf("s1=%p,cap=%d,len = %d %T \n", &s1, cap(s1), len(s1), s1)
	fmt.Printf("a=%p,cap=%d,len = %d %T\n", &a, cap(a), len(a), a)
}
