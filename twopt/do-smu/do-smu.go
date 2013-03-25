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
	"time"
)

func main() {

	// Define the input parameters
	var (
		nworkers            int
		meshsize, subsample float64
		Dfn, Rfn, cpuprof   string
		Ns, Nmu             int
		maxs                float64
		outprefix           string
	)

	// Basic flags
	flag.StringVar(&Dfn, "dfn", "", "D filename")
	flag.StringVar(&Rfn, "rfn", "", "D filename")
	flag.StringVar(&outprefix, "outprefix", "", "Output prefix")
	flag.IntVar(&Ns, "ns", 100, "Number of s bins")
	flag.IntVar(&Nmu, "nmu", 100, "Number of mu bins")
	flag.Float64Var(&maxs, "maxs", 200, "Maximum s value")
	// Tuning flags
	flag.IntVar(&nworkers, "nworkers", 1, "Number of workers")
	flag.Float64Var(&meshsize, "meshsize", 50, "Mesh size")
	flag.Float64Var(&subsample, "subsample", 1.01, "Subsampling fraction")
	flag.StringVar(&cpuprof, "cpuprofile", "", "CPU Filename")
	flag.Parse()
	if Dfn == "" {
		log.Fatal(errors.New("A data filename must be specified"))
	}
	if Rfn == "" {
		log.Fatal(errors.New("A random filename must be specified"))
	}
	if outprefix == "" {
		log.Fatal(errors.New("An output prefix must be specified"))
	}

	// Read in the particle data
	pD, err := mesh.ReadParticles(Dfn, subsample)
	boxDmin, boxDmax := pD.MinMax()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data read in from %s \n", Dfn)
	pR, err := mesh.ReadParticles(Rfn, subsample)
	boxRmin, boxRmax := pR.MinMax()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data read in from %s \n", Rfn)

	// Set up meshes	
	boxmin := boxRmin.Min(boxDmin)
	boxmax := boxRmax.Max(boxDmax)
	mD := mesh.New(pD, meshsize, boxmin, boxmax)
	mR := mesh.New(pR, meshsize, boxmin, boxmax)
	fmt.Println("Meshes created....")

	// Set up profiling if desired
	if cpuprof != "" {
		fp, err := os.Create(cpuprof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
		defer pprof.StopCPUProfile()
	}

	var hfinal twopt.PairCounter

	fmt.Println("Using nworkers=", nworkers)
	fore := twopt.NewForeman(nworkers, func() twopt.PairCounter {
		return twopt.PairCounter(twopt.NewSMuPairCounter(Ns, Nmu, maxs))
	})

	// DD
	t1 := time.Now()
	fmt.Println("Starting DD....")
	c1 := mD.LoopAll()
	for g1 := range c1 {
		c2 := mD.LoopNear(g1.I, maxs)
		for g2 := range c2 {
			switch {
			case (g1.N < g2.N):
				fore.SubmitJob(twopt.NewJob(g1, g2, 2))
			case (g1.N == g2.N):
				fore.SubmitJob(twopt.NewJob(g1, g2, 1))
			}
		}
	}
	fore.EndWork()
	hfinal = fore.Summarize()
	fmt.Printf("Time to complete DD = %s\n", time.Since(t1))
	fout, err := os.Create(outprefix + "-DD.dat")
	if err != nil {
		log.Fatal(err)
	}
	hfinal.PPrint(fout)
	fout.Close()


	_ = mR
}
