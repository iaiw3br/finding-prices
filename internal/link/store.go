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
		SELECT DISTINCT ON (iis.id) iis.id, iis.url,
                           p.price, iis.item_id, iis.store_id,
                           s.title
		FROM item_in_store iis
				 LEFT JOIN prices p ON iis.id = p.item_store_id
				 JOIN stores s on iis.store_id = s.id
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
		err = rows.Scan(&i.ItemInStore.ID, &i.ItemInStore.URL,
			&i.Price, &i.ItemInStore.ID, &i.Store.ID,
			&i.Store.Title)
		if err != nil {
			return nil, err
		}

		s.ItemInStore = i.ItemInStore
		s.Store = i.Store
		if i.Price != nil {
			s.Price.Price = *i.Price
		}

		itemsSearch = append(itemsSearch, s)
	}

	return itemsSearch, nil
}
