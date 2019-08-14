package main

import (
	"log"
	"time"
)

type Philosopher interface {
	Act()
}

type Table struct {
	Diners []Philosopher
	Forks  []*Fork

	Duration time.Duration
}

func NewTable(
	numDiners int,
	duration time.Duration,
	newDiner func(int, *Fork, *Fork) Philosopher,
) *Table {
	forks := make([]*Fork, numDiners)
	for i := 0; i < numDiners; i++ {
		forks[i] = NewFork(i)
	}

	diners := []Philosopher{}
	for i := 0; i < numDiners; i++ {
		leftFork := forks[i]
		rightFork := forks[(i+1)%numDiners]
		d := newDiner(i, leftFork, rightFork)
		diners = append(diners, d)
	}
	return &Table{
		Diners:   diners,
		Forks:    forks,
		Duration: duration,
	}
}

func (t *Table) Serve() {
	for {
		for _, d := range t.Diners {
			go d.Act()
		}

		// TODO: wait till done all steps

		<-time.Tick(t.Duration)
		log.Println()
	}
}

func main() {
	numDiners := 5
	table := NewTable(numDiners, time.Second*2, NewSimple)
	table.Serve()
}
