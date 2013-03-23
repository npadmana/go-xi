package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/npadmana/go-xi/mesh"
	"github.com/npadmana/go-xi/twopt"
	"log"
	"time"
)

func main() {
	var subsample float64
	var fn string
	flag.Float64Var(&subsample, "subfraction", 1.01, "Subsampling fraction")
	flag.StringVar(&fn, "fn", "", "Filename")
	flag.Parse()

	if fn == "" {
		log.Fatal(errors.New("A filename must be specified"))
	}

	p, err := mesh.ReadParticles(fn, subsample)
	fmt.Println("Read in Particles")
	if err != nil {
		log.Fatal(err)
	}

	t1 := time.Now()
	cc := twopt.NewSMuPairCounter(5, 5, 200)
	cc.Count(p, p, 1.0)
	dt := time.Since(t1)

	fmt.Println(dt)

}
