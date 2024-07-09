package pokedex

import "encoding/json"

type pokemonSpecies struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []name `json:"names"`
}

type name struct {
	Name     string   `json:"name"`
	Language language `json:"language"`
}

type language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func fetchPokemonSpeciesByID(id int) (*pokemonSpecies, error) {
	resp, err := speciesAPI(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var species pokemonSpecies
	if err := json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, err
	}

	return &species, nil
}

func fetchPokemonJapaneseNameByID(id int) (string, error) {
	ps, err := fetchPokemonSpeciesByID(id)
	if err != nil {
		return "", err
	}

	return getJapaneseName(ps), nil
}

func getJapaneseName(ps *pokemonSpecies) string {
	for _, name := range ps.Names {
		if name.Language.Name == "ja" {
			return name.Name
		}
	}
	return ""
}