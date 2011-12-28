package main

import (
	"fmt"
	"log"
	"time"
)

var _ = log.Printf

type Color string
type Position string

const (
	red  Color = "red"
	blue Color = "blue"

	rope  Position = "rope"
	left  Position = "left"
	right Position = "right"
)

type Baboon struct {
	id    int
	color Color
	pos   Position
}

type Rope struct {
	id int
}

func main() {
	baboons := []*Baboon{}
	for i := 0; i < 100; i++ {
		c := red
		p := left
		if i%2 == 0 {
			c = blue
			p = right
		}
		b := &Baboon{id: i, color: c, pos: p}
		baboons = append(baboons, b)
		go b.live()
	}

	ropes := []*Rope{}
	for i := 0; i < 10; i++ {
		r := &Rope{id: i}
		go r.hang()
		ropes = append(ropes, r)
	}

	fmt.Printf("Hello world %#v, %#v\n", baboons, ropes)
	select {}
}

// The baboon's lifecycle.
func (b *Baboon) live() {
	for {
		select {}
	}
}

// The riope's lifecycle.
func (r *Rope) hang() {
	tick := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			r.moveBaboons()
		}
	}
}

func (r *Rope) moveBaboons() {
	log.Printf("moving baboons on rope %d", r.id)
	
}
