package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const baseURL = "https://pokeapi.co/api/v2"

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
}

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

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {
	init := flag.Bool("init", false, "initialize the pokedex")
	flag.Parse()

	if *init {
		fmt.Println("initialize the pokedex")
		if err := fetchAllPokemonAndSaveToFile("index.csv"); err != nil {
			return err
		}
	}
	return nil
}

func fetchAllPokemonAndSaveToFile(filename string) error {
	f, err := os.Create(filename)
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
		url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", baseURL, limit, offset)
		resp, err := http.Get(url)
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

func pokemonURLToID(url string) string {
	// URLの最後のスラッシュの前にある数字をIDとして抽出
	for i := len(url) - 2; i >= 0; i-- {
		if url[i] == '/' {
			idStr := url[i+1 : len(url)-1]
			id, _ := strconv.Atoi(idStr)
			return strconv.Itoa(id)
		}
	}
	return ""
}
