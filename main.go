package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kotaoue/gopoke/pkg/pokemon"
)

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	init := flag.Bool("init", false, "initialize the pokedex")
	flag.Parse()

	// pokemon.InitializeDetails(1)

	if *init {
		fmt.Println("initialize the pokedex")
		if err := pokemon.InitializeIndex(); err != nil {
			return err
		}
	}
	return nil
}
