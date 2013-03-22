package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/twopt"
	//"github.com/npadmana/go-xi/utils"
	"log"
	"os"
	"runtime/pprof"
)

type Job struct {
	g1, g2 *mesh.GridPoint
	Scale  float64
}

type Worker struct {
	H    twopt.PairCounter
	Work chan Job
	Done chan bool
}

type Foreman struct {
	Workers    []*Worker
	LastWorker int
}

func NewForeman(n int, newpair func() twopt.PairCounter) (f *Foreman) {
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

func NewWorker(newpair func() twopt.PairCounter) (w *Worker) {
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
			w1.H.Count(job1.g1.P, job1.g2.P, job1.Scale)
		}
	}(w)
	return
}

func main() {
	var nworkers int
	var meshsize float64
	var fn, cpuprof string
	flag.IntVar(&nworkers, "nworkers", 1, "Number of workers")
	flag.Float64Var(&meshsize, "meshsize", 50, "Mesh size")
	flag.StringVar(&fn, "fn", "", "Filename")
	flag.StringVar(&cpuprof, "cpuprofile", "", "CPU Filename")
	flag.Parse()

	if fn == "" {
		log.Fatal(errors.New("A filename must be specified"))
	}

	p, err := mesh.ReadParticles(fn)
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}
	m := mesh.New(p, meshsize)
	fmt.Println("Mesh created")

	if cpuprof != "" {
		fp, err := os.Create("smu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Using nworkers=", nworkers)
	fore := NewForeman(nworkers, func() twopt.PairCounter {
		return twopt.PairCounter(twopt.NewSMuPairCounter(5, 5, 200))
	})
	c1 := m.LoopAll()
	auto := true
	for g1 := range c1 {
		c2 := m.LoopNear(g1.I, 200)
		for g2 := range c2 {
			switch {
			case !auto:
				fore.SubmitJob(Job{g1, g2, 1})
			case auto && (g1.N < g2.N):
				fore.SubmitJob(Job{g1, g2, 2})
			case auto && (g1.N == g2.N):
				fore.SubmitJob(Job{g1, g2, 1})
			}
		}
	}
	fore.EndWork()

	hfinal := fore.Workers[0].H
	for i := 1; i < len(fore.Workers); i++ {
		hfinal.Add(fore.Workers[i].H)
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Print(hfinal.Get(i, j), " ")
		}
		fmt.Println()
	}

}
