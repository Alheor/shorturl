package repository

import (
	"context"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := pg.Conn.Exec(ctx,
		"INSERT INTO "+tableName+" (correlation_id, original_url) VALUES (@correlationID, @originalURL)",
		pgx.NamedArgs{"correlationID": id, "originalURL": value},
	)

	if err != nil {
		pgError := err.(*pgconn.PgError)
		if pgError.Code == "23505" {
			return errors.New(ErrValueAlreadyExist)
		}

		panic(err)
	}

	return nil
}

func (pg *Postgres) Get(id string) (value string, error error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var originalURL string

	row := pg.Conn.QueryRow(ctx,
		"SELECT original_url FROM "+tableName+" WHERE correlation_id=@correlationID",
		pgx.NamedArgs{"correlationID": id},
	)

	err := row.Scan(&originalURL)
	if err != nil {
		if err != pgx.ErrNoRows {
			panic(err)
		}

		return ``, errors.New(ErrIDNotFound)
	}

	return originalURL, nil
}

func (pg *Postgres) Remove(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := pg.Conn.Exec(ctx,
		"DELETE FROM "+tableName+" WHERE correlation_id=@correlationID",
		pgx.NamedArgs{"correlationID": id},
	)

	if err != nil {
		panic(err)
	}
}

func (pg *Postgres) StorageIsReady() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := pg.Conn.Ping(ctx)
	return err == nil
}

func (pg *Postgres) AddBatch(in []BatchEl) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pg.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}

	for _, v := range in {
		batch.Queue("INSERT INTO "+tableName+" (correlation_id, original_url) VALUES (@correlationID, @originalURL)",
			pgx.NamedArgs{"correlationID": v.ShortURL, "originalURL": v.OriginalURL},
		)
	}

	err = tx.SendBatch(ctx, batch).Close()
	if err != nil {
		tx.Rollback(ctx)

		pgError := err.(*pgconn.PgError)

		if pgError.Code == "23505" {
			return errors.New(ErrValueAlreadyExist)
		}

		panic(err)
	}

	return tx.Commit(ctx)
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
		    correlation_id varchar(`+strconv.Itoa(randomname.ShortNameLength)+`) NOT NULL UNIQUE,
		    original_url text NOT NULL UNIQUE
		);
	`)

	if err != nil {
		panic(err)
	}
}
