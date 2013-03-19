package main

import (
	"fmt"
	"github.com/npadmana/go-xi/particle"
	"log"
)

func main() {
	parr, err := particle.NewFromXYZW("testparticlereader.dat")
	if err != nil {
		log.Fatal(err)
	}
	for _,p1 := range parr {
		fmt.Println(p1)
	}
}
