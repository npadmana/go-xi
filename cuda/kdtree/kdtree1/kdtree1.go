package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/kdtree"
	"github.com/npadmana/go-xi/cuda/particle"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	var subsample, maxs float64
	var minpart int
	var minbox float64
	var fn string
	flag.Float64Var(&subsample, "subfraction", 1.01, "Subsampling fraction")
	flag.Float64Var(&maxs, "maxs", 200, "maximum s value")
	flag.Float64Var(&minbox, "minbox", 100, "Minimum box size")
	flag.IntVar(&minpart, "minpart", 10000, "Minimum number of particles")
	flag.StringVar(&fn, "fn", "", "Filename")
	flag.Parse()

	if fn == "" {
		log.Fatal(errors.New("A filename must be specified"))
	}

	// Read in particles on the host
	p, err := particle.ReadParticles(fn, float32(subsample))
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}

	// Build tree
	t1 := time.Now()
	root := kdtree.NewNode(p, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go root.Grow(minpart, minbox, wg)
	wg.Wait()
	dt := time.Since(t1)
	fmt.Println("Time to build the tree:", dt)

}
