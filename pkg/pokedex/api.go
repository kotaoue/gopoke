package pokedex

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const graphqlURL = "https://beta.pokeapi.co/graphql/v1beta2"

type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func graphqlAPI(query string, variables map[string]interface{}) (*http.Response, error) {
	body, err := json.Marshal(graphqlRequest{
		Query:     query,
		Variables: variables,
	})
	if err != nil {
		return nil, err
	}
	return http.Post(graphqlURL, "application/json", bytes.NewReader(body))
}
