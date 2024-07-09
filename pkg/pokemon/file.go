package pokemon

import (
	"encoding/csv"
	"fmt"
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
		pokemons, err := fetchPokemonList(limit, offset)
		if err != nil {
			return err
		}

		for _, result := range pokemons.Results {
			id := urlToID(result.URL)

			if id != 0 {
				pokemon, err := fetchPokemonByID(id)
				if err != nil {
					return err
				}

				fmt.Printf("%+v\n", *pokemon)
				if err := w.Write(pokemon.toCSV()); err != nil {
					return err
				}
			}
		}

		// 次のページがある場合、オフセットを更新して次のリクエストの準備
		if pokemons.Next != "" {
			offset += limit
		} else {
			break // 次のページがない場合は終了
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
