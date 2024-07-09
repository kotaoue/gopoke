package pokemon

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pokemon struct {
	ID     int
	Name   string
	Height int
	Weight int
}

func (p Pokemon) toHeader() []string {
	return []string{
		"ID",
		"Name",
		"Height",
		"Weight",
	}
}

func (p *Pokemon) toCSV() []string {
	return []string{
		fmt.Sprintf("%d", p.ID),
		p.Name,
		fmt.Sprintf("%d", p.Height),
		fmt.Sprintf("%d", p.Weight),
	}
}

type pokemonDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

func fetchPokemonByID(id int) (*Pokemon, error) {
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
	resp, err := pokemonAPI(id)
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
