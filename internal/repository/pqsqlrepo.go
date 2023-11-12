package repository

import (
	"context"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"time"
)

const tableName = `short_URL`

// Postgres connection structure
type Postgres struct {
	Conn *pgxpool.Pool
}

func (pg *Postgres) Init() error {

	if pg.Conn != nil {
		return nil
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, config.Options.DatabaseDsn)
	if err != nil {
		panic(err)
	}

	pg.Conn = db

	createDBSchema(ctx, pg.Conn)

	return nil
}

func (pg *Postgres) Add(id string, value string) error {
	return nil
}

func (pg *Postgres) Get(id string) (value string, error error) {
	return ``, nil
}

func (pg *Postgres) Remove(id string) {

}

func (pg *Postgres) StorageIsReady() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := pg.Conn.Ping(ctx)
	if err != nil {
		return false
	}

	return true
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var tableExists bool

	row := conn.QueryRow(ctx, `SELECT true FROM pg_tables WHERE tablename = $1`, strings.ToLower(tableName))
	err := row.Scan(&tableExists)
	if err != nil {
		if err != pgx.ErrNoRows {
			panic(err)
		}

		tableExists = false
	}

	if tableExists {
		return
	}

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
		    id SERIAL NOT NULL PRIMARY KEY,
		    short varchar(`+strconv.Itoa(randomname.ShortNameLength)+`),
		    long text
		);
		CREATE UNIQUE INDEX `+tableName+`_short_long_uniq_idx ON `+tableName+` (short, long);
	`)

	if err != nil {
		panic(err)
	}
}
