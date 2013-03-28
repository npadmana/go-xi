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

func sumOfWeights(m1 *mesh.Mesh) (sum float64) {
	// Ensure that sum is really initialized, it should be, but old habits die hard
	sum = 0

	for _, p := range m1.Particles {
		sum += p.W
	}
	return
}

func doOne(m1, m2 *mesh.Mesh, outfn string, Nmu, Ns int, maxs float64, auto bool, nworkers int) {
	fore := twopt.NewForeman(nworkers, func() twopt.PairCounter {
		return twopt.PairCounter(twopt.NewSMuPairCounter(Ns, Nmu, maxs))
	})

	t1 := time.Now()
	c1 := m1.LoopAll()
	for g1 := range c1 {
		c2 := m2.LoopNear(g1.I, maxs)
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
	hfinal := fore.Summarize()
	fmt.Printf("Time to complete = %s\n", time.Since(t1))
	fout, err := os.Create(outfn)
	if err != nil {
		log.Fatal(err)
	}
	hfinal.PPrint(fout)
	fout.Close()
}

func main() {

	// Define the input parameters
	var (
		nworkers            int
		meshsize, subsample float64
		Dfn, Rfn, cpuprof   string
		Ns, Nmu             int
		maxs                float64
		outprefix           string
		help                bool
	)

	// Basic flags
	flag.BoolVar(&help, "help", false, "Prints the help message and quits")
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
	if help {
		flag.Usage()
		os.Exit(0)
	}
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

	fmt.Println("Using nworkers=", nworkers)
	// Do DD, DR, RR
	fmt.Println("Starting DD....")
	doOne(mD, mD, outprefix+"-DD.dat", Nmu, Ns, maxs, true, nworkers)
	fmt.Println("Starting DR....")
	doOne(mD, mR, outprefix+"-DR.dat", Nmu, Ns, maxs, false, nworkers)
	fmt.Println("Starting RR....")
	doOne(mR, mR, outprefix+"-RR.dat", Nmu, Ns, maxs, true, nworkers)


	// Print the auxiliary file
	fout, err := os.Create(outprefix+"-norm.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer fout.Close()
	fmt.Fprintln(fout, "# Sum of weights in input files")
	fmt.Fprintf(fout, "%s: %20.15e\n", Dfn, sumOfWeights(mD))
	fmt.Fprintf(fout, "%s: %20.15e\n", Rfn, sumOfWeights(mR))
	

}
