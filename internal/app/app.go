package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"prices/internal/customers"
	"prices/internal/customers/ipiter"
	lite_mobile "prices/internal/customers/lite-mobile"
	"prices/internal/customers/pitergsm"
	"prices/internal/customers/store78"
	"prices/internal/link"
	"prices/internal/price"
	"prices/pkg/client/postgresql"
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
	connectorRegistry := createConnector()

	for _, item := range itemsForSearch {

		doc, err := getDocument(item.ItemStore.URL)
		if err != nil {
			log.Fatal(err)
		}

		conn := connectorRegistry.Get(item.ItemStore.Store.Title)
		priceFromWebsite, err := conn.Search(doc)

		if err != nil {
			log.Fatal(err)
			//continue
		}

		if item.Price != priceFromWebsite {
			cp := price.CreatePrice{
				ItemStoreId: item.ID,
				Price:       priceFromWebsite,
				Created:     now,
			}
			err = priceService.Create(ctx, cp)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("price was been changed item id:%v\n", item.ItemStore.ItemID)
			fmt.Printf("old price: %v, new price: %v\n", item.Price, priceFromWebsite)
		}
	}
}

func createConnector() *customers.Registry[customers.GDS] {
	pg := pitergsm.NewConnector()
	lm := lite_mobile.NewConnector()
	ip := ipiter.NewConnector()
	s78 := store78.NewConnector()

	connectorRegistry := customers.GlobalRegistry()
	connectorRegistry.Add(pg, lm, ip, s78)

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
