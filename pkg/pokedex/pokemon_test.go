package pokedex

import (
	"testing"
)

func TestPokemon_toHeader(t *testing.T) {
	p := Pokemon{}
	expected := []string{"ID", "Name", "Genera", "Height", "Weight", "FlavorText"}
	result := p.toHeader()

	if len(result) != len(expected) {
		t.Errorf("Expected %d headers, got %d", len(expected), len(result))
	}

	for i, header := range expected {
		if result[i] != header {
			t.Errorf("Expected header[%d] to be %s, got %s", i, header, result[i])
		}
	}
}

func TestPokemon_toCSV(t *testing.T) {
	p := Pokemon{
		ID:         1,
		Name:       "フシギダネ",
		Genera:     "たねポケモン",
		Height:     70.0,
		Weight:     6.9,
		FlavorText: "背中に植物の種が植わっている。",
	}

	expected := []string{"1", "フシギダネ", "たねポケモン", "70.0", "6.9", "背中に植物の種が植わっている。"}
	result := p.toCSV()

	if len(result) != len(expected) {
		t.Errorf("Expected %d CSV fields, got %d", len(expected), len(result))
	}

	for i, field := range expected {
		if result[i] != field {
			t.Errorf("Expected CSV field[%d] to be %s, got %s", i, field, result[i])
		}
	}
}

func TestPokemon_toCSV_DecimalFormatting(t *testing.T) {
	p := Pokemon{
		ID:     25,
		Name:   "ピカチュウ",
		Height: 40.6,
		Weight: 6.0,
	}

	result := p.toCSV()
	
	if result[3] != "40.6" {
		t.Errorf("Expected height to be formatted as 40.6, got %s", result[3])
	}
	
	if result[4] != "6.0" {
		t.Errorf("Expected weight to be formatted as 6.0, got %s", result[4])
	}
}