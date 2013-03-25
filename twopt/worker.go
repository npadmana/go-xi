package twopt

import (
	"github.com/npadmana/go-xi/mesh"
)

type Job struct {
	g1, g2 *mesh.GridPoint
	scale  float64
}

func NewJob(g1, g2 *mesh.GridPoint, scale float64) (j Job) {
	j.g1 = g1
	j.g2 = g2
	j.scale = scale
	return
}

type Worker struct {
	H    PairCounter
	Work chan Job
	Done chan bool
}

type Foreman struct {
	Workers    []*Worker
	LastWorker int
}

func NewForeman(n int, newpair func() PairCounter) (f *Foreman) {
	f = new(Foreman)
	f.Workers = make([]*Worker, n)
	for i := range f.Workers {
		f.Workers[i] = NewWorker(newpair)
	}
	f.LastWorker = 0
	return
}

func (f *Foreman) EndWork() {
	for _, w := range f.Workers {
		close(w.Work)
	}
	for _, w := range f.Workers {
		<-w.Done
	}
}

func (f *Foreman) SubmitJob(j Job) {
	ok := false
	for !ok {
		select {
		case f.Workers[f.LastWorker].Work <- j:
			ok = true
			f.LastWorker = (f.LastWorker + 1) % len(f.Workers)
		default:
			f.LastWorker = (f.LastWorker + 1) % len(f.Workers)
		}
	}
}

func NewWorker(newpair func() PairCounter) (w *Worker) {
	w = new(Worker)
	w.H = newpair()
	w.Work = make(chan Job, 5)
	w.Done = make(chan bool)
	go func(w1 *Worker) {
		ok := true
		var job1 Job
		for {
			job1, ok = <-w1.Work
			if !ok {
				w1.Done <- true
				return
			}
			w1.H.Count(job1.g1.P, job1.g2.P, job1.scale)
		}
	}(w)
	return
}

func (f *Foreman) Summarize() PairCounter {
	hfinal := f.Workers[0].H
	for i := 1; i < len(f.Workers); i++ {
		hfinal.Add(f.Workers[i].H)
	}

	return hfinal
}

func (f *Foreman) Reset() {
	for i := 0; i < len(f.Workers); i++ {
		f.Workers[i].H.Reset()
	}
}
