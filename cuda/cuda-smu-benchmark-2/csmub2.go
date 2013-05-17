package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/mesh"
	"github.com/npadmana/go-xi/cuda/particle"
	"log"
	"os"
	"time"
)

func main() {
	var subsample, maxs float64
	var fn string
	flag.Float64Var(&subsample, "subfraction", 1.01, "Subsampling fraction")
	flag.Float64Var(&maxs, "maxs", 200, "maximum s value")
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

	// Build a mesh
	boxmin, boxmax := p.MinMax()
	m1 := mesh.New(p, 200.0, boxmin, boxmax)

	// Allocate particle data on the device 
	devP := cudalib.NewDevParticle(len(p))
	defer devP.Free()

	// Copy particle data
	t1 := time.Now()
	devP.CopyToDevice(m1.Particles)
	dt := time.Since(t1)
	fmt.Println("Time to move data onto GPU:", dt)

	smu := cudalib.NewSMuCudaPairCounter(5, 5, float32(maxs), 1.e8)
	defer smu.Free()
	t1 = time.Now()

	c1 := m1.LoopAll()
	for g1 := range c1 {
		c2 := m1.LoopNear(g1.I, float32(maxs))
		for g2 := range c2 {
			smu.Count(devP, devP, g1.Lo, g1.Hi, g2.Lo, g2.Hi, 1)
		}
	}
	//smu.Count(devP, devP, 0, len(p), 0, len(p), 1)
	smu.PullFromDevice()
	dt = time.Since(t1)
	fmt.Println("Time to pair count:", dt)
	smu.PPrint(os.Stdout)

}
