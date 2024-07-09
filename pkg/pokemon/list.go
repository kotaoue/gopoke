package pokemon

import (
	"encoding/json"
)

type pokemonList struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []pokemonListResults `json:"results"`
}

type pokemonListResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func fetchPokemonList(limit, offset int) (*pokemonList, error) {
	resp, err := pokemonListAPI(limit, offset)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ps pokemonList
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		return nil, err
	}

	return &ps, nil
}
