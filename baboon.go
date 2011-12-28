package main

import "fmt"


type Color string
type Position string

const (
	red Color = "red"
	blue Color = "blue"

	rope Position = "rope"
	left Position = "left"
	right Position = "right"
)


type Baboon struct {
	id int
	color Color
	pos Position
}


type Rope struct {
	id int
}


func main() {
	baboons := []*Baboon{}
	for i := 0; i < 100; i++ {
		c := red
		p := left
		if i % 2 == 0 {
			c = blue
			p = right
		}
		baboons = append(baboons, &Baboon{id: i, color: c, pos: p})
	}

	ropes := []*Rope{}
	for i := 0; i < 10; i++ {
		ropes = append(ropes, &Rope{id: i})
	}

	fmt.Printf("Hello world %#v, %#v\n", baboons, ropes)
}
