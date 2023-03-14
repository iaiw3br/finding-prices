package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"prices/internal/customers"
	"prices/internal/customers/ipiter"
	"prices/internal/customers/iport"
	lite_mobile "prices/internal/customers/lite-mobile"
	"prices/internal/customers/pitergsm"
	"prices/internal/customers/store78"
	"prices/internal/link"
	"prices/internal/price"
	"prices/pkg/client/postgresql"
)

func Run() {
	start := time.Now()
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
	connectorRegistry := createConnector()

	results := make(chan price.CreatePrice)

	for _, s := range itemsForSearch {
		go func(s link.Search) {
			doc, err := getDocument(s.ItemInStore.URL)
			if err != nil {
				fmt.Printf("parse url is failed: %s\n", s.ItemInStore.URL)
				return
			}

			conn := connectorRegistry.Get(s.Store.Title)
			priceFromWebsite, err := conn.Search(doc)
			if err != nil {
				fmt.Printf("get price from website failed: %s\n", s.ItemInStore.URL)
				return
			}

			if s.Price.Price != priceFromWebsite {
				fmt.Printf("send result price: %f, storeId: %d\n", priceFromWebsite, s.ItemInStore.ID)
				results <- price.CreatePrice{
					ItemStoreId: s.ItemInStore.ID,
					Price:       priceFromWebsite,
					Created:     now,
				}
			}
			close(results)
		}(s)
	}

	for result := range results {
		if _, ok := <-results; !ok {
			break
		}
		fmt.Printf("get result price: %f, storeId: %d\n", result.Price, result.ItemStoreId)
		if err = priceService.Create(ctx, result); err != nil {
			log.Fatal(err)
		}
		fmt.Println("created")
	}

	end := time.Now()
	fmt.Printf("search completed, time:%d sec\n", end.Sub(start).Milliseconds())
}

func createConnector() *customers.Registry[customers.GDS] {
	pg := pitergsm.NewConnector()
	lm := lite_mobile.NewConnector()
	ip := ipiter.NewConnector()
	s78 := store78.NewConnector()
	ipt := iport.NewConnector()

	connectorRegistry := customers.GlobalRegistry()
	connectorRegistry.Add(pg, lm, ip, s78, ipt)

	return connectorRegistry
}

func getDocument(url string) (*goquery.Document, error) {
	// Make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status is not OK")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
