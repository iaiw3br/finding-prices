package lite_mobile

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func Search(doc *goquery.Document) (float64, error) {
	var title string
	// Find the review items
	doc.Find(".price__now").Each(func(_ int, s *goquery.Selection) {
		// For each link found, get the title
		title = s.Find("span").Text()
		title = strings.ReplaceAll(title, " ", "")
	})

	price, err := strconv.ParseFloat(title, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}
