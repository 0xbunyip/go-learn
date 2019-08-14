package main

import "log"

type Simple struct {
	Id    int
	State int

	LeftFork  *Fork
	RightFork *Fork
}

func NewSimple(id int, leftFork *Fork, rightFork *Fork) Philosopher {
	return &Simple{
		Id:        id,
		State:     0,
		LeftFork:  leftFork,
		RightFork: rightFork,
	}
}

func (s *Simple) Act() {
	if s.State == 0 {
		if got := TryFork(s.LeftFork, s.Id); got > 0 {
			s.State = 1
			log.Printf("Philosopher %d acquired left fork %d", s.Id, s.LeftFork.Id)
		} else {
			log.Printf("Philosopher %d failed to acquire left fork %d", s.Id, s.LeftFork.Id)
		}
	} else if s.State == 1 {
		if got := TryFork(s.RightFork, s.Id); got > 0 {
			s.State = 2
			log.Printf("Philosopher %d acquired right fork %d", s.Id, s.RightFork.Id)
		} else {
			log.Printf("Philosopher %d failed to acquire right fork %d", s.Id, s.RightFork.Id)
		}
	} else if s.State == 2 {
		s.Eat() // TODO: check if both forks are acquired
		s.State = 3
	} else if s.State == 3 {
		ReturnFork(s.RightFork, s.Id)
		s.State = 4
		log.Printf("Philosopher %d released right fork %d", s.Id, s.RightFork.Id)
	} else if s.State == 4 {
		ReturnFork(s.LeftFork, s.Id)
		s.State = 0
		log.Printf("Philosopher %d released left fork %d", s.Id, s.LeftFork.Id)
	}
}

func (s *Simple) Eat() {
	log.Printf("Philosopher %d is eating", s.Id)
}
