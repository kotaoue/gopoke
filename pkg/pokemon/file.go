package pokemon

import (
	"os"
	"path/filepath"
)

const (
	indexFile = "pokedex/index.csv"
)

func dirAndFileCreate(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(name)
}
