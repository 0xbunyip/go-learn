package main

import (
	"log"
	"time"
)

type Philosopher interface {
	Act(chan<- int)
}

type Table struct {
	Diners []Philosopher
	Forks  []*Fork

	Duration  time.Duration
	Spaghetti chan int
	Step      int
}

func NewTable(
	numDiners int,
	duration time.Duration,
	newDiner func(int, *Fork, *Fork, chan int) Philosopher,
) *Table {
	forks := make([]*Fork, numDiners)
	for i := 0; i < numDiners; i++ {
		forks[i] = NewFork(i)
	}

	spaghetti := make(chan int, 1000)

	diners := []Philosopher{}
	for i := 0; i < numDiners; i++ {
		leftFork := forks[i]
		rightFork := forks[(i+1)%numDiners]
		d := newDiner(i, leftFork, rightFork, spaghetti)
		diners = append(diners, d)
	}
	t := &Table{
		Diners:    diners,
		Forks:     forks,
		Duration:  duration,
		Spaghetti: spaghetti,
		Step:      0,
	}
	go t.Refill()
	return t
}

func (t *Table) Serve() {
	done := make(chan int, len(t.Diners))
	for {
		for _, d := range t.Diners {
			go d.Act(done)
		}

		for i := 0; i < len(t.Diners); i++ {
			<-done
		}

		time.Sleep(t.Duration)
		t.Step++
		log.Println()
	}
}

func (t *Table) Refill() {
	count := make([]int, len(t.Diners))
	for {
		id := <-t.Spaghetti
		count[id]++
		log.Printf("Spaghetti at %d: %v", t.Step, count)
	}
}

func main() {
	numDiners := 5
	table := NewTable(numDiners, time.Second*2, NewParity)
	table.Serve()
}
