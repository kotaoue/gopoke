package pokemon

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Pokemon struct {
	ID     int
	Name   string
	Height int
	Weight int
}

type pokemonDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

func FetchPokemonByID(id int) (*Pokemon, error) {
	pd, err := fetchPokemonDetailByID(id)
	if err != nil {
		return nil, err
	}

	name, err := fetchPokemonJapaneseNameByID(id)
	if err != nil {
		return nil, err
	}

	return &Pokemon{
		ID:     pd.ID,
		Name:   name,
		Height: pd.Height,
		Weight: pd.Weight,
	}, nil
}

func fetchPokemonDetailByID(id int) (*pokemonDetail, error) {
	resp, err := fetchPokemonByID(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var pd pokemonDetail
	if err := json.NewDecoder(resp.Body).Decode(&pd); err != nil {
		return nil, err
	}

	return &pd, nil
}

func InitializeDetails(id int) error {
	f, err := dirAndFileCreate(detailFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"id", "name", "weight", "height"}); err != nil {
		return err
	}

	// id分だけループする
	resp, err := fetchPokemonByID(id)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var pokemon pokemonDetail
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return err
	}

	fmt.Printf("%+v\n", pokemon)
	record := []string{
		strconv.Itoa(pokemon.ID),
		pokemon.Name,
		strconv.Itoa(pokemon.Height),
		strconv.Itoa(pokemon.Weight),
	}
	if err := w.Write(record); err != nil {
		return err
	}

	return nil
}
