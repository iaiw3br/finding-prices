package link

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store interface {
	ItemsForSearch(ctx context.Context) ([]Search, error)
}

type repository struct {
	client *pgxpool.Pool
}

func NewStore(client *pgxpool.Pool) Store {
	return &repository{
		client: client,
	}
}

func (r repository) ItemsForSearch(ctx context.Context) ([]Search, error) {
	sql := `
		SELECT DISTINCT ON (lis.id) lis.id, lis.url, p.price, lis.item_id, lis.store_id
		FROM link_items_stores lis
		LEFT JOIN prices p ON lis.id = p.item_store_id
		ORDER BY id, created DESC;
	`

	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemsSearch []Search
	for rows.Next() {
		var i SearchNil
		var s Search
		err = rows.Scan(&i.ID, &i.ItemStore.URL, &i.Price, &i.ItemStore.ItemID, &i.ItemStore.StoreID)
		if err != nil {
			return nil, err
		}

		s.ID = i.ID
		s.ItemStore.URL = i.ItemStore.URL
		s.ItemStore.StoreID = i.ItemStore.StoreID
		if i.Price != nil {
			s.Price = *i.Price
		}

		itemsSearch = append(itemsSearch, s)
	}

	return itemsSearch, nil
}
