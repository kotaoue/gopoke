package pokemon

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
)

type pokemonListResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []pokemonDetail `json:"results"`
}

type pokemonListItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func InitializeIndex() error {
	f, err := dirAndFileCreate(indexFile)
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
		resp, err := fetchPokemons(limit, offset)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// レスポンスデータの取得
		var result struct {
			Count    int               `json:"count"`
			Next     string            `json:"next"`
			Previous string            `json:"previous"`
			Results  []pokemonListItem `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		for _, pokemon := range result.Results {
			record := []string{
				urlToID(pokemon.URL),
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
