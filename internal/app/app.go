package app

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"prices/internal/link"
	"prices/internal/price"
	"prices/pkg/client/postgresql"
	"strconv"
	"strings"
	"time"
)

func Run() {
	ctx := context.Background()

	client, err := postgresql.New(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	priceStore := price.NewStore(client)
	priceService := price.NewService(priceStore)
	linkItemStore := link.NewStore(client)

	linkService := link.NewService(linkItemStore)
	itemsForSearch, err := linkService.FindForSearch(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("elements for search: %d\n", len(itemsForSearch))

	now := time.Now()
	for _, item := range itemsForSearch {
		priceFromWebsite := getPriceFromWebsite(item.ItemStore.URL)

		if item.Price != priceFromWebsite {
			cp := price.CreatePrice{
				ItemStoreId: item.ID,
				Price:       priceFromWebsite,
				Created:     now,
			}
			err = priceService.Create(ctx, cp)
			if err != nil {
				log.Fatal(err)
				// need continue
			}
			fmt.Printf("price was been changed item id:%v\n", item.ItemStore.ItemID)
			fmt.Printf("old price: %v, new price: %v\n", item.Price, priceFromWebsite)
		}
	}
}

func getPriceFromWebsite(url string) float64 {
	// Make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var title string
	// Find the review items
	doc.Find(".price__now").Each(func(_ int, s *goquery.Selection) {
		// For each link found, get the title
		title = s.Find("span").Text()
		title = strings.ReplaceAll(title, " ", "")
	})

	price, err := strconv.ParseFloat(title, 64)
	if err != nil {
		log.Fatal(err)
	}

	return price
}
