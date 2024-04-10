package main

import (
	"log"
	"naboj.org/letter/pkg"
)

func main() {
	log.Printf("Letter v.%v\n", pkg.VERSION)

	r := pkg.NewRouter()
	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}
