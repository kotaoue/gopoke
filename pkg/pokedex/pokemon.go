package pokedex

import "fmt"

type Pokemon struct {
	ID         int
	Name       string
	Genera     string
	Height     float64
	Weight     float64
	FlavorText string
}

func (p Pokemon) toHeader() []string {
	return []string{
		"ID",
		"Name",
		"Genera",
		"Height",
		"Weight",
		"FlavorText",
	}
}

func (p *Pokemon) toCSV() []string {
	return []string{
		fmt.Sprintf("%d", p.ID),
		p.Name,
		p.Genera,
		fmt.Sprintf("%.1f", p.Height),
		fmt.Sprintf("%.1f", p.Weight),
		p.FlavorText,
	}
}

func PrintPokemon(p Pokemon) error {
	fmt.Printf("No.%04d\n", p.ID)
	fmt.Printf("%s\n", p.Name)
	fmt.Printf("分類: %s\n", p.Genera)
	fmt.Printf("高さ: %.1fcm\t重さ: %.1fkg\n", p.Height, p.Weight)
	fmt.Println()
	fmt.Println(p.FlavorText)
	return nil
}
