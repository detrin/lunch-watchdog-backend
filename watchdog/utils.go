package watchdog

import (
	"regexp"
	"strconv"
	"strings"
)

func parsePrice(priceText string) (int, error) {
	reg := regexp.MustCompile(`[0-9]+`)
	priceText = strings.TrimSpace(reg.FindString(priceText))
	price, err := strconv.Atoi(priceText)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func parseAllergens(description string) []string {
	allergensRegex := regexp.MustCompile(`\|([0-9,]+)\|`)
	matches := allergensRegex.FindStringSubmatch(description)
	if len(matches) > 1 {
		return strings.Split(matches[1], ",")
	}
	return nil
}

func normalizeSpace(input string) string {
	withSpaces := strings.ReplaceAll(input, "\u00A0", " ")
	spacePattern := regexp.MustCompile(`\s+`)
	collapsedSpaces := spacePattern.ReplaceAllString(withSpaces, " ")

	return collapsedSpaces
}
