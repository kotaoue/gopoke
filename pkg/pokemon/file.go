package pokemon

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func InitializePokedex() error {
	f, err := dirAndFileCreate("pokedex/pokedex.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(Pokemon{}.toHeader()); err != nil {
		return err
	}

	limit := 100
	offset := 0

	for {
		ps, err := fetchPokemonList(limit, offset)
		if err != nil {
			return err
		}

		for _, result := range ps.Results {
			id := urlToID(result.URL)

			log.Printf("Saving Pokemon ID: %d name: %s\n", id, result.Name)

			// IDs 10000 and above are special forms and do not have Japanese names
			if id == 0 || id >= 10000 {
				continue
			}

			pokemon, err := fetchPokemonByID(id)
			if err != nil {
				return fmt.Errorf("failed to fetch pokemon by id: %w", err)
			}

			if err := w.Write(pokemon.toCSV()); err != nil {
				return fmt.Errorf("failed to write csv: %w", err)
			}
		}

		// Update the offset to prepare for the next request if there has a next page.
		if ps.Next != "" {
			offset += limit
		} else {
			break
		}
	}

	return nil
}

func dirAndFileCreate(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(name)
}
