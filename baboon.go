package main

import "fmt"


type Color string

const (
	RED Color = "red"
	BLUE Color = "blue"
)


type Baboon struct {
	number int
	color Color
}


func main() {
	b := Baboon{2, RED}
	fmt.Printf("Hello world %#v\n", b)
}
