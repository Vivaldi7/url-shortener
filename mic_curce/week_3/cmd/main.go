package main

import "fmt"

func main() {
	res := Div(1.0, 2.0)
	fmt.Println(res)
}

func Div(a float64, b float64) float64 {
	return a / b
}
