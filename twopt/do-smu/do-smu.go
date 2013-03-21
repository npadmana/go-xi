package main

import (
	"fmt"
	//"github.com/npadmana/go-xi/particle"
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/twopt"
	"github.com/npadmana/go-xi/utils"
	"log"
)

type Job struct {
	g1, g2 *mesh.GridPoint
}

type Worker struct {
	H    *utils.UniformHistogram
	Work chan Job
	Done chan bool
}

type Foreman struct {
	Workers []*Worker
}

func NewForeman(n int) (f *Foreman) {
	f = new(Foreman)
	f.Workers = make([]*Worker, n)
	for i := range f.Workers {
		f.Workers[i] = NewWorker()
	}
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
	i := 0
	for !ok {
		select {
		case f.Workers[i].Work <- j:
			ok = true
		default:
			i = (i + 1) % len(f.Workers)
		}
	}
}

func NewWorker() (w *Worker) {
	w = new(Worker)
	w.H = utils.NewUniform([]int{5, 5}, []float64{0, 0}, []float64{200, 1})
	w.Work = make(chan Job)
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
			twopt.PairCounter(w1.H, job1.g1.P, job1.g2.P, twopt.SMu)
		}
	}(w)
	return
}

func main() {
	p, err := mesh.ReadParticles("test_N.dat")
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}
	m := mesh.New(p, 50.0)
	fmt.Println("Mesh created")

	fore := NewForeman(8)
	c1 := m.LoopAll()
	for g1 := range c1 {
		c2 := m.LoopAll()
		for g2 := range c2 {
			//_, _ = g1, g2
			//twopt.PairCounter(w1.H, g1.P, g2.P, twopt.SMu)
			fore.SubmitJob(Job{g1, g2})
		}
	}
	fore.EndWork()

	hfinal := fore.Workers[0].H
	for i := 1; i < 8; i++ {
		hfinal.AddHist(fore.Workers[i].H)
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(hfinal.Get(i, j), " ")
		}
		fmt.Println()
	}

}
