package link

import "time"

type itemInStore struct {
	ID     int
	URL    string
	ItemID int
}

type SearchNil struct {
	ItemInStore itemInStore
	Store       store
	Price       *float64
	Item        item
}

type Search struct {
	ItemInStore itemInStore
	Store       store
	Price       price
	Item        item
}

type price struct {
	Price   float64
	Created time.Time
}

type store struct {
	ID    int
	Title string
}

type item struct {
	ID    int
	Title string
}
