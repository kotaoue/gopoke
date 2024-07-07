package pokemon

import (
	"fmt"
	"net/http"
	"strconv"
)

const baseURL = "https://pokeapi.co/api/v2"

func pokemonAPI(limit, offset int) (*http.Response, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d&offset=%d", baseURL, limit, offset)
	return http.Get(url)
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
