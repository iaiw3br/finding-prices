package price

import "time"

type CreatePrice struct {
	ItemStoreId int
	Price       float64
	Created     time.Time
}
