package ipiter

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func Search(doc *goquery.Document) (float64, error) {
	var title string
	// Find the review items
	doc.Find(".pricebox").Each(func(_ int, s *goquery.Selection) {
		// For each link found, get the title
		title = s.Find(".price").Text()
		title = strings.ReplaceAll(title, " ", "")
		title = strings.ReplaceAll(title, "\n", "")
		title = strings.ReplaceAll(title, "₽", "")
	})

	price, err := strconv.ParseFloat(title, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
