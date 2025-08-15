package pokedex

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test_pokedex.db")

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := createPokedexTable(db); err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(testDBPath)
	}

	return db, cleanup
}

func insertTestPokemon(t *testing.T, db *sql.DB, pokemon Pokemon) {
	query := "INSERT INTO pokemons (id, name, genera, height, weight, flavor_text) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, pokemon.ID, pokemon.Name, pokemon.Genera, pokemon.Height, pokemon.Weight, pokemon.FlavorText)
	if err != nil {
		t.Fatalf("Failed to insert test pokemon: %v", err)
	}
}

func TestCreatePokedexTable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	rows, err := db.Query("PRAGMA table_info(pokemons)")
	if err != nil {
		t.Fatalf("Failed to query table info: %v", err)
	}
	defer rows.Close()

	columnCount := 0
	for rows.Next() {
		columnCount++
	}

	expectedColumns := 6
	if columnCount != expectedColumns {
		t.Errorf("Expected %d columns, got %d", expectedColumns, columnCount)
	}
}

func TestExistsPokemons(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	exists, err := existsPokemons(db)
	if err != nil {
		t.Fatalf("Failed to check if pokemons table exists: %v", err)
	}

	if !exists {
		t.Error("Expected pokemons table to exist")
	}
}

func TestCountPokemons(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	count, err := countPokemons(db)
	if err != nil {
		t.Fatalf("Failed to count pokemons: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 pokemons, got %d", count)
	}

	testPokemon := Pokemon{
		ID:         1,
		Name:       "フシギダネ",
		Genera:     "たねポケモン",
		Height:     70.0,
		Weight:     6.9,
		FlavorText: "背中に植物の種が植わっている。",
	}
	insertTestPokemon(t, db, testPokemon)

	count, err = countPokemons(db)
	if err != nil {
		t.Fatalf("Failed to count pokemons after insert: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 pokemon, got %d", count)
	}
}

func TestSelectPokemon(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	testPokemon := Pokemon{
		ID:         25,
		Name:       "ピカチュウ",
		Genera:     "ねずみポケモン",
		Height:     40.0,
		Weight:     6.0,
		FlavorText: "電気を貯める袋を頬に持つ。",
	}
	insertTestPokemon(t, db, testPokemon)

	result, err := SelectPokemon(db, 25)
	if err != nil {
		t.Fatalf("Failed to select pokemon: %v", err)
	}

	if result.ID != testPokemon.ID {
		t.Errorf("Expected ID %d, got %d", testPokemon.ID, result.ID)
	}
	if result.Name != testPokemon.Name {
		t.Errorf("Expected name %s, got %s", testPokemon.Name, result.Name)
	}
	if result.Genera != testPokemon.Genera {
		t.Errorf("Expected genera %s, got %s", testPokemon.Genera, result.Genera)
	}
}

func TestSelectPokemon_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := SelectPokemon(db, 999)
	if err == nil {
		t.Error("Expected error when selecting non-existent pokemon")
	}

	expectedError := "pokemon not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRandomPokemon(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	testPokemon := Pokemon{
		ID:         1,
		Name:       "フシギダネ",
		Genera:     "たねポケモン",
		Height:     70.0,
		Weight:     6.9,
		FlavorText: "背中に植物の種が植わっている。",
	}
	insertTestPokemon(t, db, testPokemon)

	result, err := RandomPokemon(db)
	if err != nil {
		t.Fatalf("Failed to get random pokemon: %v", err)
	}

	if result.ID != testPokemon.ID {
		t.Errorf("Expected ID %d, got %d", testPokemon.ID, result.ID)
	}
}

func TestRandomPokemon_EmptyTable(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := RandomPokemon(db)
	if err == nil {
		t.Error("Expected error when getting random pokemon from empty table")
	}

	expectedError := "pokemon not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSelectPokemons_ByName(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	testPokemons := []Pokemon{
		{ID: 1, Name: "フシギダネ", Height: 70.0, Weight: 6.9},
		{ID: 25, Name: "ピカチュウ", Height: 40.0, Weight: 6.0},
		{ID: 150, Name: "ミュウツー", Height: 200.0, Weight: 122.0},
	}

	for _, p := range testPokemons {
		insertTestPokemon(t, db, p)
	}

	searchCondition := SearchCondition{
		Name:  "%ピカ%",
		Limit: 10,
	}

	results, err := SelectPokemons(db, searchCondition)
	if err != nil {
		t.Fatalf("Failed to select pokemons by name: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "ピカチュウ" {
		t.Errorf("Expected name ピカチュウ, got %s", results[0].Name)
	}
}

func TestSelectPokemons_ByHeight(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	testPokemons := []Pokemon{
		{ID: 1, Name: "フシギダネ", Height: 70.0, Weight: 6.9},
		{ID: 25, Name: "ピカチュウ", Height: 40.0, Weight: 6.0},
	}

	for _, p := range testPokemons {
		insertTestPokemon(t, db, p)
	}

	searchCondition := SearchCondition{
		Height: 50.0,
		Limit:  10,
	}

	results, err := SelectPokemons(db, searchCondition)
	if err != nil {
		t.Fatalf("Failed to select pokemons by height: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestSelectQuery_InvalidCondition(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	searchCondition := SearchCondition{
		Height: -1,
		Weight: -1,
		Name:   "",
		Limit:  10,
	}

	_, err := selectQuery(db, searchCondition)
	if err == nil {
		t.Error("Expected error for invalid search condition")
	}

	expectedError := "Please set a value greater than 0 for either weight or height"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
