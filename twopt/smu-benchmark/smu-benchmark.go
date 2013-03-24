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
	var subsample, maxs float64
	var fn, cpuprof string
	flag.Float64Var(&subsample, "subfraction", 1.01, "Subsampling fraction")
	flag.Float64Var(&maxs, "maxs", 200, "maximum s value")
	flag.StringVar(&fn, "fn", "", "Filename")
	flag.StringVar(&cpuprof, "cpuprofile", "", "CPU Filename")
	flag.Parse()

	if fn == "" {
		log.Fatal(errors.New("A filename must be specified"))
	}

	if cpuprof != "" {
		fp, err := os.Create("smu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
		defer pprof.StopCPUProfile()
	}

	p, err := mesh.ReadParticles(fn, subsample)
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}

	t1 := time.Now()
	cc := twopt.NewSMuPairCounter(5, 5, maxs)
	cc.Count(p, p, 1.0)
	dt := time.Since(t1)

	fmt.Println(dt)

}
