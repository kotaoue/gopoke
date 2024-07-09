package pokemon

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PokemonDetail struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

type PokemonSpecies struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []Name `json:"names"`
}

func FetchPokemonByID(id int) (*PokemonDetail, error) {
	resp, err := fetchPokemonByID(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var pokemon PokemonDetail
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", pokemon)
	return &pokemon, nil
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

	var pokemon PokemonDetail
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

	species, err := getPokemonSpecies(id)
	if err != nil {
		return err
	}

	japaneseName := getJapaneseName(species)
	fmt.Printf("Japanese name for Pokémon with ID %d: %s\n", id, japaneseName)

	return nil
}

func getPokemonSpecies(id int) (*PokemonSpecies, error) {
	resp, err := fetchSpeciesByID(id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var species PokemonSpecies
	if err := json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, err
	}

	return &species, nil
}

func getJapaneseName(species *PokemonSpecies) string {
	for _, name := range species.Names {
		if name.Language.Name == "ja" {
			return name.Name
		}
	}
	return ""
}
