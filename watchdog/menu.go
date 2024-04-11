package watchdog

import "time"

type Menu struct {
	Name      string     `json:"name"`
	Date      time.Time  `json:"date"`
	MenuItems []MenuItem `json:"menu_items"`
}

type MenuItem struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}
