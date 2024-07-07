package pokemon

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
)

type PokemonListResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Pokemon `json:"results"`
}

type PokemonListItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func FetchAllPokemonAndSaveToFile(filename string) error {
	f, err := dirAndFileCreate(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"id", "name", "url"}); err != nil {
		return err
	}

	limit := 100
	offset := 0

	for {
		resp, err := pokemonAPI(limit, offset)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// レスポンスデータの取得
		var result struct {
			Count    int               `json:"count"`
			Next     string            `json:"next"`
			Previous string            `json:"previous"`
			Results  []PokemonListItem `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		for _, pokemon := range result.Results {
			record := []string{
				pokemonURLToID(pokemon.URL),
				pokemon.Name,
				pokemon.URL,
			}
			fmt.Println(record)
			if err := w.Write(record); err != nil {
				return err
			}
		}

		// 次のページがある場合、オフセットを更新して次のリクエストの準備
		if result.Next != "" {
			offset += limit
		} else {
			break // 次のページがない場合は終了
		}
	}

	return nil
}
