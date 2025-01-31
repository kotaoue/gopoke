package pokedex

import (
	"encoding/json"
)

type pokemonSpecies struct {
	ID                int          `json:"id"`
	Name              string       `json:"name"`
	Names             []name       `json:"names"`
	FlavorTextEntries []flavorText `json:"flavor_text_entries"`
	Genera            []genus      `json:"genera"`
}

type flavorText struct {
	FlavorText string `json:"flavor_text"`
	Language   language
}
type name struct {
	Name     string   `json:"name"`
	Language language `json:"language"`
}
type genus struct {
	Genus    string   `json:"genus"`
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

func getJapaneseName(ps *pokemonSpecies) string {
	for _, name := range ps.Names {
		if name.Language.Name == "ja" {
			return name.Name
		}
	}
	return ""
}

func getJapaneseFlavorText(ps *pokemonSpecies) string {
	for _, ft := range ps.FlavorTextEntries {
		if ft.Language.Name == "ja" {
			return ft.FlavorText
		}
	}
	return ""
}

func getJapaneseGenes(ps *pokemonSpecies) string {
	for _, g := range ps.Genera {
		if g.Language.Name == "ja" {
			return g.Genus
		}
	}
	return ""
}
