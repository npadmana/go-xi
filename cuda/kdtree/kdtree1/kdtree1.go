package main

import (
	"errors"
	"flag"
	"fmt"
	//"github.com/npadmana/go-xi/cuda/cudalib"
	"github.com/npadmana/go-xi/cuda/kdtree"
	"github.com/npadmana/go-xi/cuda/particle"
	"log"
	//"os"
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
	root := kdtree.NewNode(p, 0, len(p), 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go root.Grow(minpart, float32(minbox), &wg)
	wg.Wait()
	dt := time.Since(t1)
	fmt.Println("Time to build the tree:", dt)

	fmt.Printf("The root has %d particles \n", root.Npart)

	done := make(chan bool)
	out := make(chan kdtree.NodeList)
	npart := 0
	go func() {
		for n1 := range out {
			//fmt.Printf("This is node %d with %d particles....\n", n1[0].Id, n1[0].Npart)
			npart += n1[0].Npart
		}
		done <- true
	}()
	wg.Add(1)
	t1 = time.Now()
	go kdtree.TreeMap(root, func(n1 kdtree.NodeList) kdtree.TreeDecision { return kdtree.CONTINUE }, out, &wg)
	wg.Wait()
	close(out)
	<-done
	dt = time.Since(t1)
	fmt.Println("Time to walk the tree and count the nodes:", dt)

	fmt.Printf("The leaves have %d particles \n", npart)

	npairs := int64(0)
	nexec := 0
	out = make(chan kdtree.NodeList)
	go func() {
		for n1 := range out {
			npairs += int64(n1[0].Npart * n1[1].Npart)
			nexec++
		}
		done <- true
	}()
	wg.Add(1)
	rr := kdtree.RInterval{0, float32(maxs)}
	t1 = time.Now()
	go kdtree.DualTreeMap(root, root, func(n1 kdtree.NodeList) kdtree.TreeDecision { return rr.DualNodeTest(n1) }, out, &wg)
	wg.Wait()
	close(out)
	<-done
	dt = time.Since(t1)
	fmt.Println("Time to walk the trees, collecting pairs:", dt)
	fmt.Printf("Number of pairs of nodes considered : %d\n", nexec)
	fmt.Printf("Number of pairs : %d \n", npairs)
	fmt.Printf("Fractional number of pairs : %f \n", float64(npairs)/(float64(root.Npart)*float64(root.Npart)))

}
