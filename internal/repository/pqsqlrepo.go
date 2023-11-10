package repository

import (
	"context"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
	"time"
)

var (
	pgOnce sync.Once
)

type Postgres struct {
	Conn *pgxpool.Pool
}

func (sn *Postgres) Init() error {

	pgOnce.Do(func() {
		db, err := pgxpool.New(context.Background(), config.Options.DatabaseDsn)
		if err != nil {
			panic(err)
		}

		sn.Conn = db
	})

	return nil
}

func (sn *Postgres) Add(id string, value string) error {
	return nil
}

func (sn *Postgres) Get(id string) (value string, error error) {
	return ``, nil
}

func (sn *Postgres) Remove(id string) {

}

func (sn *Postgres) StorageIsReady() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := sn.Conn.Ping(ctx)
	if err != nil {
		return false
	}

	return true
}
