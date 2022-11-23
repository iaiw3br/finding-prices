package main

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type CreateRow struct {
	Title   string
	Price   int
	Code    string
	Updated time.Time
}

func main() {
	codes := []string{
		"13263",
		"13262",
		"13267",
		"13266",
		"23299",
		"23298",
	}

	ctx := context.Background()
	client, err := connectDB(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range codes {
		url := fmt.Sprintf("https://pitergsm.ru/catalog/tablets-and-laptops/mac/macbook-pro/macbook-pro-14-2021/%s/", c)
		doSomething(ctx, url, c, client)
	}

}

func doSomething(ctx context.Context, url, code string, client *pgxpool.Pool) {
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
		// For each item found, get the title
		title = s.Find("span").Text()
		title = strings.ReplaceAll(title, " ", "")
	})

	price, err := strconv.Atoi(title)
	if err != nil {
		log.Fatal(err)
	}

	cr := CreateRow{
		Title:   "macbook-pro-14-2021",
		Price:   price,
		Code:    code,
		Updated: time.Now(),
	}

	priceDB := FindByCode(ctx, cr, client)
	if cr.Price != priceDB {
		Create(ctx, cr, client)
	}

}

func connectDB(ctx context.Context) (*pgxpool.Pool, error) {
	username := "postgres"
	password := "postgres"
	tableName := "products"
	//username := os.Getenv("USERNAME")
	//password := os.Getenv("PASSWORD")
	//tableName := os.Getenv("TABLE_NAME")
	connString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", username, password, tableName)

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func FindByCode(ctx context.Context, cr CreateRow, client *pgxpool.Pool) int {
	query := `
		SELECT price
		FROM products
		WHERE code = $1
		ORDER BY updated_at desc
		LIMIT 1;
	`
	var price int
	_ = client.QueryRow(ctx, query, cr.Code).Scan(&price)

	return price
}

func Create(ctx context.Context, cr CreateRow, client *pgxpool.Pool) error {
	query := `
		INSERT INTO products (title, price, code, updated_at)
		VALUES ($1, $2, $3, $4);
	`
	_, err := client.Exec(ctx, query, cr.Title, cr.Price, cr.Code, cr.Updated)
	if err != nil {
		return err
	}
	return nil
}
