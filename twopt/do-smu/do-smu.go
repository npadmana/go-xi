package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/twopt"
	"log"
	"os"
	"runtime/pprof"
)


func main() {
	var nworkers int
	var meshsize, subsample float64
	var fn, cpuprof string
	flag.IntVar(&nworkers, "nworkers", 1, "Number of workers")
	flag.Float64Var(&meshsize, "meshsize", 50, "Mesh size")
	flag.Float64Var(&subsample, "subfraction", 1.01, "Subsampling fraction")
	flag.StringVar(&fn, "fn", "", "Filename")
	flag.StringVar(&cpuprof, "cpuprofile", "", "CPU Filename")
	flag.Parse()

	if fn == "" {
		log.Fatal(errors.New("A filename must be specified"))
	}

	p, err := mesh.ReadParticles(fn, subsample)
	boxmin, boxmax := p.MinMax()
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}
	m := mesh.New(p, meshsize, boxmin, boxmax)
	fmt.Println("Mesh created")

	if cpuprof != "" {
		fp, err := os.Create(cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
		defer pprof.StopCPUProfile()
	}

	fmt.Println("Using nworkers=", nworkers)
	fore := twopt.NewForeman(nworkers, func() twopt.PairCounter {
		return twopt.PairCounter(twopt.NewSMuPairCounter(5, 5, 200))
	})
	c1 := m.LoopAll()
	auto := true
	for g1 := range c1 {
		c2 := m.LoopNear(g1.I, 200)
		for g2 := range c2 {
			switch {
			case !auto:
				fore.SubmitJob(twopt.NewJob(g1, g2, 1))
			case auto && (g1.N < g2.N):
				fore.SubmitJob(twopt.NewJob(g1, g2, 2))
			case auto && (g1.N == g2.N):
				fore.SubmitJob(twopt.NewJob(g1, g2, 1))
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
