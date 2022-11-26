package price

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store interface {
	Create(ctx context.Context, cp CreatePrice) error
}

type repository struct {
	client *pgxpool.Pool
}

func NewStore(client *pgxpool.Pool) Store {
	return &repository{
		client: client,
	}
}

func (r repository) Create(ctx context.Context, cp CreatePrice) error {
	sql := `
		INSERT INTO prices (created, item_store_id, price) 
		VALUES ($1, $2, $3);
		`

	_, err := r.client.Exec(ctx, sql, cp.Created, cp.ItemStoreId, cp.Price)
	if err != nil {
		return err
	}
	return nil
}
