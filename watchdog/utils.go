package watchdog

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	gtranslate "github.com/gilang-as/google-translate"
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

func translate(text string) (string, error) {
	value := gtranslate.Translate{
		Text: text,
		//From: "cz",
		To: "en",
	}
	translated, err := gtranslate.Translator(value)
	if err != nil {
		return "", err
	} else {
		prettyJSON, err := json.MarshalIndent(translated, "", "\t")
		if err != nil {
			panic(err)
		}
		// fmt.Println(string(prettyJSON))

		var jsonResult struct {
			Text string `json:"text"`
		}
		json.Unmarshal(prettyJSON, &jsonResult)
		return jsonResult.Text, nil
	}
}
