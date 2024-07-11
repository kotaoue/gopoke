package pokedex

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	pokedexCSV = "pokedex/pokedex.csv"
	pokedexDB  = "pokedex/pokedex.db"
)

type SearchCondition struct {
	Height float64
	Weight float64
	Name   string
	Limit  int
}

func Initialize() error {
	if _, err := os.Stat(pokedexCSV); err == nil {
		fmt.Printf("%s already exists. Skipping creation.\n", pokedexCSV)
	} else if !os.IsNotExist(err) {
		return err
	} else {
		if err := createPokedexCSV(); err != nil {
			return err
		}
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	b, err := existsPokemons(db)
	if err != nil {
		return err
	}

	if b {
		fmt.Printf("%s already exists. Skipping creation.\n", pokedexDB)
	} else {
		if err := createPokedexTable(db); err != nil {
			return err
		}

		if err := importCSVToSQLite(db); err != nil {
			return err
		}

		i, err := countPokemons(db)
		if err != nil {
			return err
		}
		fmt.Printf("pokemons count: %d\n", i)
	}

	return nil
}

func createPokedexCSV() error {
	f, err := dirAndFileCreate(pokedexCSV)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(Pokemon{}.toHeader()); err != nil {
		return err
	}

	limit := 100
	offset := 0

	for {
		ps, err := fetchPokemonList(limit, offset)
		if err != nil {
			return err
		}

		for _, result := range ps.Results {
			id := urlToID(result.URL)

			log.Printf("Saving Pokemon ID: %d name: %s\n", id, result.Name)

			// IDs 10000 and above are special forms and do not have Japanese names
			if id == 0 || id >= 10000 {
				continue
			}

			pokemon, err := fetchPokemonByID(id)
			if err != nil {
				return fmt.Errorf("failed to fetch pokemon by id: %w", err)
			}

			if err := w.Write(pokemon.toCSV()); err != nil {
				return fmt.Errorf("failed to write csv: %w", err)
			}
		}

		// Update the offset to prepare for the next request if there has a next page.
		if ps.Next != "" {
			offset += limit
		} else {
			break
		}
	}

	return nil
}

func OpenDB() (*sql.DB, error) {
	return sql.Open("sqlite3", pokedexDB)
}

func dirAndFileCreate(name string) (*os.File, error) {
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(name)
}

func createPokedexTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pokemons (
		id INTEGER PRIMARY KEY,
		name TEXT,
		height REAL,
		weight REAL
	);`
	_, err := db.Exec(query)
	return err
}

func importCSVToSQLite(db *sql.DB) error {
	file, err := os.Open(pokedexCSV)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO pokemons (id, name, height, weight) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		id := record[0]
		name := record[1]
		height := record[2]
		weight := record[3]

		_, err = stmt.Exec(id, name, height, weight)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func existsPokemons(db *sql.DB) (bool, error) {
	rows, err := db.Query("PRAGMA table_info(pokemons)")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func countPokemons(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM pokemons").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func SelectPokemons(db *sql.DB, search SearchCondition) ([]Pokemon, error) {
	rows, err := selectQuery(db, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pokemons []Pokemon
	for rows.Next() {
		var p Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Height, &p.Weight); err != nil {
			return nil, err
		}
		pokemons = append(pokemons, p)
	}
	return pokemons, nil
}

func selectQuery(db *sql.DB, sc SearchCondition) (*sql.Rows, error) {
	if sc.Name != "" {
		query := `
		SELECT id, name, height, weight
		FROM pokemons
		WHERE name LIKE ?
		ORDER BY name
		LIMIT ?`

		return db.Query(query, sc.Name, sc.Limit)
	}

	if sc.Height >= 0 && sc.Weight >= 0 {
		query := `
		SELECT id, name, height, weight
		FROM pokemons
		ORDER BY (ABS(height - ?) + ABS(weight - ?)) / 2
		LIMIT ?`

		return db.Query(query, sc.Height, sc.Weight, sc.Limit)
	}

	if sc.Height >= 0 {
		query := `
		SELECT id, name, height, weight
		FROM pokemons
		ORDER BY ABS(height - ?), weight
		LIMIT ?`

		return db.Query(query, sc.Height, sc.Limit)
	}

	if sc.Weight >= 0 {
		query := `
		SELECT id, name, height, weight
		FROM pokemons
		ORDER BY ABS(weight - ?), height
		LIMIT ?`

		return db.Query(query, sc.Weight, sc.Limit)
	}

	return nil, errors.New("Please set a value greater than 0 for either weight or height")
}
