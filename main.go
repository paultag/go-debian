package main

import (
	"./dependency"
	"fmt"
	"log"
)

func main() {
	dep, err := dependency.Parse("foo, bar [amd64 i386] | baz (>= 1.0)")
	if err != nil {
		log.Fatalf("Oh no! %s", err)
	}
	for _, relation := range dep.Relations {
		for _, possi := range relation.Possibilities {
			fmt.Printf("   -> %s\n", possi.Name)
			for _, arch := range possi.Arches {
				fmt.Printf("      %s\n", arch.Name)
			}
		}
		fmt.Printf(".\n")
	}
	fmt.Printf("\n")
}
