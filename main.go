package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kotaoue/gopoke/pkg/pokedex"
)

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	init := flag.Bool("init", false, "initialize the pokedex")
	height := flag.Float64("height", -1, "height of the pokemon")
	weight := flag.Float64("weight", -1, "weight of the pokemon")
	limit := flag.Int("limit", 10, "limit of the pokemons")
	flag.Parse()

	if *init {
		fmt.Println("initialize the pokedex")
		if err := pokedex.Initialize(); err != nil {
			return err
		}
	}

	db, err := pokedex.OpenDB()
	if err != nil {
		return err
	}

	sc := pokedex.SearchCondition{
		Height: *height,
		Weight: *weight,
		Limit:  *limit,
	}
	ps, err := pokedex.SelectPokemons(db, sc)
	if err != nil {
		return err
	}

	for _, p := range ps {
		fmt.Printf("%s %.1fcm %.1fkg\n", p.Name, p.Height, p.Weight)
	}

	return nil
}
