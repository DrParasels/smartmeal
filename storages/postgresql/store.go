package postgresql

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/smartmeal/storages/postgresql/schema/sqlc"
)

type Storage struct {
	q *sqlc.Queries
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{q: sqlc.New(pool)}
}
