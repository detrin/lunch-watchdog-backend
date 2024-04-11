package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/detrin/lunch-watchdog-backend/watchdog"
)

func main() {
	var menus []watchdog.Menu
	menuKolkovna, err := watchdog.ScrapeMenuKolkovna()
	if err != nil {
		log.Fatalf("Error scraping menu: %v", err)
	}
	// fmt.Printf("Menu for %s - %s\n", menuKolkovna.Name, menuKolkovna.Date)
	// for i, item := range menuKolkovna.MenuItems {
	// 	fmt.Printf("Item %d: %#v\n", i+1, item)
	// }
	menus = append(menus, *menuKolkovna)

	menuMerkur, err := watchdog.ScrapeMenuMerkur()
	if err != nil {
		log.Fatalf("Error scraping menu: %v", err)
	}
	// fmt.Printf("Menu for %s - %s\n", menuMerkur.Name, menuMerkur.Date)
	// for i, item := range menuMerkur.MenuItems {
	// 	fmt.Printf("Item %d: %#v\n", i+1, item)
	// }
	menus = append(menus, *menuMerkur)

	jsonData, err := json.MarshalIndent(menus, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling menu to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
