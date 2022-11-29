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
                           s.title, i.title 
		FROM item_in_store iis
				 LEFT JOIN prices p ON iis.id = p.item_store_id
				 JOIN stores s on iis.store_id = s.id
				 JOIN items i on iis.item_id = i.id
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
		err = rows.Scan(&i.ItemInStore.ID, &i.ItemInStore.URL,
			&i.Price, &i.ItemInStore.ItemID, &i.Store.ID,
			&i.Store.Title, &i.Item.Title)
		if err != nil {
			return nil, err
		}

		itemsSearch = append(itemsSearch, convertToSearch(i))
	}

	return itemsSearch, nil
}

func convertToSearch(s SearchNil) Search {
	search := Search{
		ItemInStore: s.ItemInStore,
		Store:       s.Store,
		Item:        s.Item,
	}

	if s.Price != nil {
		search.Price.Price = *s.Price
	}

	return search
}
