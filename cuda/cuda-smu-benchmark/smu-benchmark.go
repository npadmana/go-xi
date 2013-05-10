package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/mesh"
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
	p, err := mesh.ReadParticles(fn, float32(subsample))
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}

	// Allocate particle data on the device 
	devP := cudalib.NewDevParticle(len(p))
	defer devP.Free()

	// Copy particle data
	t1 := time.Now()
	devP.CopyToDevice(p)
	dt := time.Since(t1)
	fmt.Println("Time to move data onto GPU:", dt)

	smu := cudalib.NewSMuCudaPairCounter(5, 5, float32(maxs), 1.e8)
	defer smu.Free()
	t1 = time.Now()
	smu.Count(devP, devP, 0, len(p), 0, len(p), 1)
	smu.PullFromDevice()
	dt = time.Since(t1)
	fmt.Println("Time to pair count:", dt)
	smu.PPrint(os.Stdout)

}
