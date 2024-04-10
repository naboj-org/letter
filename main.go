package main

import (
	"log"
	"naboj.org/letter/pkg"
)

func main() {
	r := pkg.NewRouter()
	err := r.Run()
	if err != nil {
		log.Panic(err)
	}
}
