package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(ctx context.Context) (*pgxpool.Pool, error) {
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
