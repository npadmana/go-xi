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
	
	parr.Normalize()
	
	for _,p1 := range parr.Data {
		fmt.Println(p1)
	}
	
	fmt.Println("BoxDim=", parr.BoxDim)
	fmt.Println("Origin=", parr.Origin)
}
