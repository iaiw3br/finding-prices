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
		SELECT lis.url, p.price, lis.id
		FROM (
			SELECT price, item_store_id,
				   row_number() OVER (PARTITION BY item_store_id ORDER BY created DESC ) AS rn
			FROM prices
			 ) p
		JOIN link_items_stores lis ON lis.id = p.item_store_id
		WHERE rn = 1;
	`

	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemsSearch []Search
	for rows.Next() {
		var i Search
		err = rows.Scan(&i.ItemStore.URL, &i.Price, &i.ID)
		if err != nil {
			return nil, err
		}

		itemsSearch = append(itemsSearch, i)
	}

	return itemsSearch, nil
}
