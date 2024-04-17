package watchdog

import "time"

type Menu struct {
	Name      string     `json:"name"`
	Date      time.Time  `json:"date"`
	MenuItems []MenuItem `json:"menu_items"`
}

type MenuItem struct {
	Description   string `json:"description"`
	DescriptionEN string `json:"description_en"`
	Price         int    `json:"price"`
}

func (mi *MenuItem) TranslateEN() error {
	translation, err := translate(mi.Description)
	if err != nil {
		return err
	}
	mi.DescriptionEN = translation
	return nil
}

func (menu *Menu) TranslateEN() error {
	for i := range menu.MenuItems {
		err := menu.MenuItems[i].TranslateEN()
		if err != nil {
			return err
		}
	}
	return nil
}
