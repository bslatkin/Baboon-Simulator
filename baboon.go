package main

import (
	"fmt"
	"log"
	"math/rand"
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

	numRopes   = 1
	numBaboons = 10
	ropeLength = 5 // Max Baboons on a Rope
)

type Baboon struct {
	// Immutable
	id    int
	color Color
	posc  chan Position // Notification from Rope that we've moved

	// Owned by Baboon's event loop:
	pos Position
}

type Rope struct {
	// Immutable.
	id                  int
	fromLeft, fromRight chan *Baboon // to move onto the rope

	// Owned by Rope's event loop:
	c        chan *Baboon // channel capacity is how many baboons can fit
	lastButt Color        // color of last dude accepted onto rope; invalid if len(c) == 0
	towards  Position     // towards left or right; invalid if len(c) == 0
}

var (
	ropes []*Rope
)

func main() {
	baboons := []*Baboon{}
	for i := 0; i < numBaboons; i++ {
		c := red
		p := left
		if i%2 == 0 {
			c = blue
			p = right
		}
		b := &Baboon{
			id:    i,
			color: c,
			pos:   p,
			posc:  make(chan Position, 1),
		}
		baboons = append(baboons, b)
		go b.live()
	}

	for i := 0; i < numRopes; i++ {
		r := &Rope{
			id:        i,
			fromRight: make(chan *Baboon),
			fromLeft:  make(chan *Baboon),
			c:         make(chan *Baboon, ropeLength),
		}
		go r.hang()
		ropes = append(ropes, r)
	}

	fmt.Printf("Hello world %#v, %#v\n", baboons, ropes)
	select {}
}

// The baboon's lifecycle.
func (b *Baboon) live() {
	for {
		if b.pos == rope {
			b.pos = <-b.posc
		} else {
			b.poop()
			// Now pick a rope to get on
			newRope := ropes[rand.Intn(len(ropes))]
			e := b.entrance(newRope)
			select {
			case e <- b:
				// Does not necessarily mean the Baboon has crossed the gorge, could be that the rope rejected the Baboon
				b.pos = rope
			case <-time.After(1000 * time.Millisecond):
				// Couldn't get a rope
			}
		}
	}
}

func (b *Baboon) entrance(r *Rope) chan<- *Baboon {
	if b.pos == right {
		return r.fromRight
	}
	return r.fromLeft
}

func (b *Baboon) poop() {
	log.Printf("%s pooping on %s side", b, b.pos)
	wait := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(wait)
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
			r.acceptBaboon(b, left)
		case b := <-fr:
			r.acceptBaboon(b, right)
		}
	}
}

func (r *Rope) acceptBaboon(b *Baboon, fromPos Position) {
	if len(r.c) > 0 && b.color != r.lastButt {
		// Send the Baboon back to the side it came from as if it never actually got on the rope.
		log.Printf("%s: %s rejected from %s due to butt", r, b, fromPos)
		b.posc <- fromPos
		return
	}
	log.Printf("%s: %s entered from %s", r, b, fromPos)
	r.towards = fromPos
	r.lastButt = b.color
	r.c <- b // can't block; cap verified
}

func (r *Rope) moveBaboons() {
	select {
	case b := <-r.c:
		log.Printf("%s: %s is now on %s side", r, b, r.towards)
		b.posc <- r.towards
	default:
		// Nothing to do.
	}
}

func (r *Rope) String() string {
	return fmt.Sprintf("Rope-%d", r.id)
}
