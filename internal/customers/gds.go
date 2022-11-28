package customers

import "github.com/PuerkitoBio/goquery"

type GDS interface {
	Search(doc *goquery.Document) (float64, error)
	Named
}
