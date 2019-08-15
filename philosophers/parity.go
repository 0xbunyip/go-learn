package main

import "log"

type Parity struct {
	Id    int
	State int

	LeftFork  *Fork
	RightFork *Fork
	Spaghetti chan int
}

func NewParity(id int, leftFork *Fork, rightFork *Fork, spaghetti chan int) Philosopher {
	return &Parity{
		Id:        id,
		State:     0,
		LeftFork:  leftFork,
		RightFork: rightFork,
		Spaghetti: spaghetti,
	}
}

func (p *Parity) Act(done chan<- int) {
	if p.State == 0 {
		fork := p.LeftFork
		if p.Id%2 == 0 {
			fork = p.RightFork
		}
		if got := TryFork(fork, p.Id); got > 0 {
			p.State = 1
			log.Printf("Philosopher %d acquired fork %d", p.Id, fork.Id)
		} else {
			log.Printf("Philosopher %d failed to acquire fork %d", p.Id, fork.Id)
		}
	} else if p.State == 1 {
		fork := p.RightFork
		if p.Id%2 == 0 {
			fork = p.LeftFork
		}
		if got := TryFork(fork, p.Id); got > 0 {
			p.State = 2
			log.Printf("Philosopher %d acquired fork %d", p.Id, fork.Id)
		} else {
			log.Printf("Philosopher %d failed to acquire fork %d", p.Id, fork.Id)
		}
	} else if p.State == 2 {
		p.Eat() // TODO: check if both forks are acquired
		p.State = 3
	} else if p.State == 3 {
		ReturnFork(p.RightFork, p.Id)
		p.State = 4
		log.Printf("Philosopher %d released right fork %d", p.Id, p.RightFork.Id)
	} else if p.State == 4 {
		ReturnFork(p.LeftFork, p.Id)
		p.State = 0
		log.Printf("Philosopher %d released left fork %d", p.Id, p.LeftFork.Id)
	}
	done <- 1
}

func (p *Parity) Eat() {
	log.Printf("Philosopher %d is eating", p.Id)
	p.Spaghetti <- p.Id
}
