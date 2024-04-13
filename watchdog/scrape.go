package watchdog

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeMenuKolkovna() (*Menu, error) {
	res, err := http.Get("https://dock.kolkovna.cz/")
	if err != nil {
		return nil, fmt.Errorf("error getting URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML document: %w", err)
	}

	var menu Menu
	var menuItems []MenuItem

	doc.Find(".op-menu-day.active").Each(func(i int, s *goquery.Selection) {
		date := s.AttrOr("data-date", "")
		date_parsed, err := time.Parse("2006-01-02", date)
		if err == nil {
			menu.Date = date_parsed
		}

		items := s.Find(".food-list-daily").Contents()
		items.Each(func(index int, sel *goquery.Selection) {
			text := sel.Text()
			if sel.Is("span") && sel.HasClass("price") {
				price, _ := parsePrice(text)
				menuItems = append(menuItems, MenuItem{Price: price})
			}
		})
		var descriptions []string
		menu_html, err := s.Find(".food-list-daily").Html()
		if err != nil {
			return
		}
		for _, row := range strings.Split(menu_html, "\n") {
			row = strings.TrimSpace(row)
			if len(row) == 0 || row[0] == '<' {
				continue
			}
			if strings.Contains(row, "<br/>") {
				new_row := strings.Split(row, "<br/>")
				if len(new_row) > 0 {
					row = new_row[0]
				}
			}
			descriptions = append(descriptions, row)
		}
		for mi, description := range descriptions {
			description = strings.Split(description, "|")[0]
			menuItems[mi].Description = description
		}
	})
	menu.Name = "Kolkovna"
	menu.MenuItems = menuItems

	return &menu, nil
}

func ScrapeMenuMerkur() (*Menu, error) {
	res, err := http.Get("http://www.restauracemerkur.cz/poledni-nabidka")
	if err != nil {
		return nil, fmt.Errorf("error getting URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error loading HTML document: %w", err)
	}

	var menu Menu
	var menuItems []MenuItem

	header := doc.Find("#main > div.post > div.post_body > h2 > strong > em > strong > em > strong > em").First()
	dateRgx := regexp.MustCompile(`\d+\.\d+\.\d+`)

	if date := dateRgx.FindString(header.Text()); date != "" {
		date_parsed, err := time.Parse("1.2.2006", date)
		if err == nil {
			menu.Date = date_parsed
		}
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		tag_description := ""
		s.Find("td:nth-child(2) > h3 > em").Each(func(i int, em *goquery.Selection) {
			tag_description = tag_description + em.Text()
		})
		s.Find("td:nth-child(2) > h3 > strong > em").Each(func(i int, em *goquery.Selection) {
			tag_description = tag_description + em.Text()
		})
		fmt.Println(tag_description)
		withRegularSpaces := normalizeSpace(tag_description)
		description := strings.TrimSpace(withRegularSpaces)
		tag_price := s.Find("td:nth-child(3)").First().Text()
		price, _ := parsePrice(tag_price)
		if price > 0 {
			menuItems = append(menuItems, MenuItem{Description: description, Price: price})
		}

	})
	menu.Name = "Merkur"
	menu.MenuItems = menuItems

	return &menu, nil
}
