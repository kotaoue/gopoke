package pokedex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const pokemonBatchQuery = `
query GetPokemonBatch($limit: Int!, $offset: Int!) {
  pokemon_v2_pokemon(
    limit: $limit,
    offset: $offset,
    order_by: {id: asc},
    where: {id: {_lt: 10000}}
  ) {
    id
    height
    weight
    pokemon_v2_pokemonspecy {
		pokemon_v2_pokemonspeciesnames(where: {language_id: {_eq: 1}}) {
        	name
			genus
		}
		pokemon_v2_pokemonspeciesflavortexts(where: {language_id: {_eq: 1}}, limit: 1) {
			flavor_text
		}
	}
  }
}
`

type graphqlResponse struct {
	Data struct {
		Pokemon []gqlPokemon `json:"pokemon_v2_pokemon"`
	} `json:"data"`
	Errors []gqlError `json:"errors"`
}

type gqlError struct {
	Message string `json:"message"`
}

type gqlPokemon struct {
	ID      int         `json:"id"`
	Height  int         `json:"height"`
	Weight  int         `json:"weight"`
	Species *gqlSpecies `json:"pokemon_v2_pokemonspecy"`
}

type gqlSpecies struct {
	Names       []gqlName       `json:"pokemon_v2_pokemonspeciesnames"`
	FlavorTexts []gqlFlavorText `json:"pokemon_v2_pokemonspeciesflavortexts"`
}

type gqlName struct {
	Name  string `json:"name"`
	Genus string `json:"genus"`
}

type gqlFlavorText struct {
	FlavorText string `json:"flavor_text"`
}

func fetchPokemonBatch(limit, offset int) ([]Pokemon, error) {
	resp, err := graphqlAPI(pokemonBatchQuery, map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var result graphqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		messages := make([]string, 0, len(result.Errors))
		for _, e := range result.Errors {
			if e.Message != "" {
				messages = append(messages, e.Message)
			}
		}
		if len(messages) == 0 {
			return nil, fmt.Errorf("graphql returned errors")
		}
		return nil, fmt.Errorf("graphql returned errors: %s", strings.Join(messages, "; "))
	}

	pokemons := make([]Pokemon, 0, len(result.Data.Pokemon))
	for _, gp := range result.Data.Pokemon {
		p := Pokemon{
			ID:     gp.ID,
			Height: float64(gp.Height) * 10.0, // convert to cm
			Weight: float64(gp.Weight) / 10.0, // convert to kg
		}
		if gp.Species != nil {
			if len(gp.Species.Names) > 0 {
				p.Name = gp.Species.Names[0].Name
				p.Genera = gp.Species.Names[0].Genus
			}
			if len(gp.Species.FlavorTexts) > 0 {
				p.FlavorText = gp.Species.FlavorTexts[0].FlavorText
			}
		}
		pokemons = append(pokemons, p)
	}

	if len(pokemons) == 0 && offset == 0 {
		return nil, fmt.Errorf("pokemon batch is empty at offset=0")
	}

	return pokemons, nil
}
