package main

import (
	"bytes"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/kotaoue/gopoke/pkg/pokedex"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestMain(t *testing.T) (*sql.DB, func()) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test_pokedex.db")

	originalPokedexDB := "pokedex/pokedex.db"

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	err = os.MkdirAll("pokedex", 0755)
	if err != nil {
		t.Fatalf("Failed to create pokedex directory: %v", err)
	}

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS pokemons (
		id INTEGER PRIMARY KEY,
		name TEXT,
		genera TEXT,
		height REAL,
		weight REAL,
		flavor_text TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	testPokemon := `INSERT INTO pokemons (id, name, genera, height, weight, flavor_text) VALUES
		(1, 'フシギダネ', 'たねポケモン', 70.0, 6.9, '背中に植物の種が植わっている。'),
		(25, 'ピカチュウ', 'ねずみポケモン', 40.0, 6.0, '電気を貯める袋を頬に持つ。')`
	_, err = db.Exec(testPokemon)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	err = os.Rename(testDBPath, originalPokedexDB)
	if err != nil {
		t.Fatalf("Failed to move test database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Chdir(oldWd)
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func captureOutput(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestPrintRandomPokemon(t *testing.T) {
	db, cleanup := setupTestMain(t)
	defer cleanup()

	output, err := captureOutput(func() error {
		return printRandomPokemon(db)
	})

	if err != nil {
		t.Fatalf("printRandomPokemon failed: %v", err)
	}

	if output == "" {
		t.Error("Expected output, got empty string")
	}

	if !bytes.Contains([]byte(output), []byte("No.")) {
		t.Error("Expected output to contain Pokemon number format")
	}
}

func TestPrintPokemonByID(t *testing.T) {
	db, cleanup := setupTestMain(t)
	defer cleanup()

	output, err := captureOutput(func() error {
		return printPokemonByID(db, 1)
	})

	if err != nil {
		t.Fatalf("printPokemonByID failed: %v", err)
	}

	if output == "" {
		t.Error("Expected output, got empty string")
	}

	if !bytes.Contains([]byte(output), []byte("No.0001")) {
		t.Error("Expected output to contain 'No.0001'")
	}

	if !bytes.Contains([]byte(output), []byte("フシギダネ")) {
		t.Error("Expected output to contain 'フシギダネ'")
	}
}

func TestPrintPokemonByID_NotFound(t *testing.T) {
	db, cleanup := setupTestMain(t)
	defer cleanup()

	_, err := captureOutput(func() error {
		return printPokemonByID(db, 999)
	})

	if err == nil {
		t.Error("Expected error when printing non-existent pokemon")
	}

	expectedError := "pokemon not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPrintPokemons(t *testing.T) {
	db, cleanup := setupTestMain(t)
	defer cleanup()

	searchCondition := pokedex.SearchCondition{
		Height: 50.0,
		Weight: -1,
		Name:   "",
		Limit:  10,
	}

	output, err := captureOutput(func() error {
		return printPokemons(db, searchCondition)
	})

	if err != nil {
		t.Fatalf("printPokemons failed: %v", err)
	}

	if output == "" {
		t.Error("Expected output, got empty string")
	}

	if !bytes.Contains([]byte(output), []byte("cm")) {
		t.Error("Expected output to contain height in cm")
	}

	if !bytes.Contains([]byte(output), []byte("kg")) {
		t.Error("Expected output to contain weight in kg")
	}
}
