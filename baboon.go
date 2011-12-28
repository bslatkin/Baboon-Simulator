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

type ropeQuery struct {
	res chan ropeStatus
}

// ropeStatus is a snapshot of the rope's status when asking
type ropeStatus struct {
	free     bool  // if no baboons on rope
	occupied Color // if !free
	towards  Position
}

type Rope struct {
	// Immutable.
	id                  int
	fromLeft, fromRight chan *Baboon // to move onto the rope
	qc                  chan ropeQuery

	// Owned by Rope's event loop:
	c        chan *Baboon // channel capacity is how many baboons can fit
	lastButt Color        // color of last dude accepted onto rope; invalid if len(c) == 0
	towards  Position     // towards left or right; invalid if len(c) == 0
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

func (b *Baboon) String() string {
	return fmt.Sprintf("Baboon-%s-%d", b.color, b.id)
}

// The riope's lifecycle.
func (r *Rope) hang() {
	tick := time.NewTicker(100 * time.Millisecond)
	for {
		nRope := len(r.c) // number of baboons on the rope
		fl, fr := r.fromLeft, r.fromRight
		switch {
		case nRope == cap(r.c):
			// Rope is too full to accept new baboons.
			fr, fl = nil, nil
		case nRope > 0:
			// If any baboons on the rope, only accept
			// from the correct direction of travel.
			switch r.towards {
			case right:
				fr = nil
			case left:
				fl = nil
			}
		}
		select {
		case <-tick.C:
			r.moveBaboons()
		case b := <-fl:
			r.towards = right
			r.c <- b // can't block; cap verified
		case b := <-fr:
			r.towards = left
			r.c <- b // can't block; cap verified
		}
	}
}

func (r *Rope) moveBaboons() {
	select {
	case b := <-r.c:
		log.Printf("%s moved %s to %s", r, b, r.towards)
		b.pos = r.towards // TODO: don't mess with its state; send baboon a message that it's been moved.
	default:
		// Nothing to do.
		log.Printf("%s idle", r)
	}
}

func (r *Rope) String() string {
	return fmt.Sprintf("Rope-%d", r.id)
}
