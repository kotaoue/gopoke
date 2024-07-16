package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pokemon struct {
	ID         int
	Name       string
	Genera     string
	Height     float64
	Weight     float64
	FlavorText string
}

func (p Pokemon) toHeader() []string {
	return []string{
		"ID",
		"Name",
		"Genera",
		"Height",
		"Weight",
		"FlavorText",
	}
}

func (p *Pokemon) toCSV() []string {
	return []string{
		fmt.Sprintf("%d", p.ID),
		p.Name,
		p.Genera,
		fmt.Sprintf("%.1f", p.Height),
		fmt.Sprintf("%.1f", p.Weight),
		p.FlavorText,
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

	ps, err := fetchPokemonSpeciesByID(id)
	if err != nil {
		return nil, err
	}

	return &Pokemon{
		ID:         pd.ID,
		Name:       getJapaneseName(ps),
		Genera:     getJapaneseGenes(ps),
		FlavorText: getJapaneseFlavorText(ps),
		Height:     float64(pd.Height) * 10.0, // convert to cm
		Weight:     float64(pd.Weight) / 10.0, // convert to kg
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

func PrintPokemonByID(id int) error {
	p, err := fetchPokemonByID(id)
	if err != nil {
		return err
	}

	ps, err := fetchPokemonSpeciesByID(id)
	if err != nil {
		return err
	}

	fmt.Printf("No.%04d\n", id)
	fmt.Printf("%s\n", getJapaneseName(ps))
	fmt.Printf("分類: %s\n", getJapaneseGenes(ps))
	fmt.Printf("高さ: %.1fcm\t重さ: %.1fkg\n", p.Height, p.Weight)
	fmt.Println()
	fmt.Println(getJapaneseFlavorText(ps))

	return nil
}
