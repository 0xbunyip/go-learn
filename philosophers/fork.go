package main

type Fork struct {
	Id int

	Requests   chan ForkRequest
	AcquiredBy int
}

func NewFork(id int) *Fork {
	f := &Fork{
		Id:         id,
		Requests:   make(chan ForkRequest, 1000),
		AcquiredBy: -1,
	}
	go f.Ready()
	return f
}

type ForkRequest struct {
	Requester int
	Action    int

	Result chan int
}

func (f *Fork) Ready() {
	for {
		r := <-f.Requests
		f.Process(r)
	}
}

func (f *Fork) Process(r ForkRequest) {
	if r.Action == 0 { // Pick up
		result := 0
		if f.AcquiredBy < 0 {
			result = 1
			f.AcquiredBy = r.Requester
		}
		r.Result <- result // TODO: prevent lock here

	} else if r.Action == 1 { // Put down
		result := 0
		if f.AcquiredBy >= 0 && f.AcquiredBy == r.Requester {
			result = 1
			f.AcquiredBy = -1
		}
		r.Result <- result // TODO: prevent lock here
	}
}

func TryFork(f *Fork, requester int) int {
	return makeForkRequest(f, requester, 0)
}

func ReturnFork(f *Fork, requester int) int {
	return makeForkRequest(f, requester, 1)
}

func makeForkRequest(f *Fork, requester, action int) int {
	res := make(chan int)
	r := ForkRequest{
		Requester: requester, // TODO: prevent id stealing
		Action:    action,
		Result:    res,
	}
	f.Requests <- r
	return <-res
}
