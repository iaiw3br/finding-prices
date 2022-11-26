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
		SELECT DISTINCT ON (lis.id) lis.id, lis.url, p.price
		FROM link_items_stores lis
		JOIN prices p ON lis.id = p.item_store_id
		ORDER BY id, created DESC;
	`

	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itemsSearch []Search
	for rows.Next() {
		var i Search
		err = rows.Scan(&i.ID, &i.ItemStore.URL, &i.Price)
		if err != nil {
			return nil, err
		}

		itemsSearch = append(itemsSearch, i)
	}

	return itemsSearch, nil
}
