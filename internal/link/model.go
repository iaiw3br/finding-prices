package link

type ItemStore struct {
	ItemID  int    `json:"item_id"`
	StoreID int    `json:"store_id"`
	URL     string `json:"url"`
}

type SearchNil struct {
	ID        int
	Price     *float64
	ItemStore ItemStore
}

type Search struct {
	ID        int
	Price     float64
	ItemStore ItemStore
}
