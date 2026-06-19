package main

import "fmt"

func foo(src []int) {
	src = append(src, 5)
	fmt.Println(src)
}

func main() {
	arc := []int{1, 2, 3}
	var src = arc[:1]

	foo(src)

	fmt.Println(src)
	fmt.Println(arc)
}
