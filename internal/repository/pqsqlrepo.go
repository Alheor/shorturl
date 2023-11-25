package repository

import (
	"context"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
)

const tableName = `short_URL`

// Postgres connection structure
type Postgres struct {
	Conn *pgxpool.Pool
}

func (pg *Postgres) Init(ctx context.Context) error {

	if pg.Conn != nil {
		return nil
	}

	db, err := pgxpool.New(ctx, config.Options.DatabaseDsn)
	if err != nil {
		panic(err)
	}

	pg.Conn = db

	createDBSchema(ctx, pg.Conn)

	return nil
}

func (pg *Postgres) Add(ctx context.Context, id string, value string) error {

	_, err := pg.Conn.Exec(ctx,
		"INSERT INTO "+tableName+" (short_key, original_url) VALUES (@shortKey, @originalURL)",
		pgx.NamedArgs{"shortKey": id, "originalURL": value},
	)

	if err == nil {
		return nil
	}

	var myErr *pgconn.PgError
	if errors.As(err, &myErr) && myErr.Code == pgerrcode.UniqueViolation {

		uniqByOriginalURL := strings.Contains(myErr.Detail, `original_url`)
		uniqByShortKey := strings.Contains(myErr.Detail, `short_key`)

		if !uniqByShortKey && !uniqByOriginalURL {
			return myErr
		}

		if uniqByShortKey {
			return NewUniqueError(id, myErr)
		}

		if uniqByOriginalURL {
			row := pg.Conn.QueryRow(ctx,
				"SELECT short_key FROM "+tableName+" WHERE original_url=@originalUrl",
				pgx.NamedArgs{"originalUrl": value},
			)

			var shortKey string
			err := row.Scan(&shortKey)
			if err != nil {
				return myErr
			}

			return NewUniqueError(shortKey, myErr)
		}
	}

	return err
}

func (pg *Postgres) Get(ctx context.Context, id string) (value string, error error) {

	var originalURL string

	row := pg.Conn.QueryRow(ctx,
		"SELECT original_url FROM "+tableName+" WHERE short_key=@shortKey",
		pgx.NamedArgs{"shortKey": id},
	)

	err := row.Scan(&originalURL)
	if err != nil {
		return ``, errors.New(ErrIDNotFound)
	}

	return originalURL, nil
}

func (pg *Postgres) Remove(ctx context.Context, id string) {

	_, err := pg.Conn.Exec(ctx,
		"DELETE FROM "+tableName+" WHERE short_key=@shortKey",
		pgx.NamedArgs{"shortKey": id},
	)

	if err != nil {
		panic(err)
	}
}

func (pg *Postgres) IsReady(ctx context.Context) bool {

	err := pg.Conn.Ping(ctx)
	return err == nil
}

func (pg *Postgres) AddBatch(ctx context.Context, in []BatchEl) error {

	tx, err := pg.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}

	for _, v := range in {
		batch.Queue("INSERT INTO "+tableName+" (short_key, original_url) VALUES (@shortKey, @originalURL)",
			pgx.NamedArgs{"shortKey": v.ShortURL, "originalURL": v.OriginalURL},
		)
	}

	err = tx.SendBatch(ctx, batch).Close()
	if err != nil {
		tx.Rollback(ctx)

		var myErr *pgconn.PgError
		if errors.As(err, &myErr) && myErr.Code == pgerrcode.UniqueViolation {
			return errors.New(ErrValueAlreadyExist)
		}

		return err
	}

	return tx.Commit(ctx)
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) {

	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
		    id SERIAL NOT NULL PRIMARY KEY,
		    short_key varchar(`+strconv.Itoa(randomname.ShortNameLength)+`) NOT NULL UNIQUE,
		    original_url text NOT NULL UNIQUE
		);
	`)

	if err != nil {
		panic(err)
	}
}
