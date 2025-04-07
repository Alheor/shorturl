package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const tableName = `short_url`

// PostgresRepo connection structure
type PostgresRepo struct {
	Conn *pgxpool.Pool
}

// Add Добавить URL
func (pg *PostgresRepo) Add(ctx context.Context, name string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	hash := urlhasher.GetShortNameGenerator().Generate()

	_, err := pg.Conn.Exec(ctx,
		"INSERT INTO "+tableName+" (short_key, original_url) VALUES (@shortKey, @originalURL)",
		pgx.NamedArgs{"shortKey": hash, "originalURL": name},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := pg.Conn.QueryRow(ctx,
				"SELECT short_key FROM "+tableName+" WHERE original_url=@originalUrl",
				pgx.NamedArgs{"originalUrl": name},
			)

			var shortKey string
			err = row.Scan(&shortKey)
			if err != nil {
				return ``, err
			}

			return ``, &models.UniqueErr{Err: pgErr, ShortKey: shortKey}
		}

		return ``, err
	}

	return hash, nil
}

// AddBatch Добавить URL пачкой
func (pg *PostgresRepo) AddBatch(ctx context.Context, list *[]models.BatchEl) error {

	tx, err := pg.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	var entries [][]any
	for _, v := range *list {
		entries = append(entries, []any{v.ShortURL, v.OriginalURL})
	}

	_, err = pg.Conn.CopyFrom(
		ctx,
		pgx.Identifier{tableName},
		[]string{"short_key", "original_url"},
		pgx.CopyFromRows(entries),
	)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// GetByShortName получить URL по короткому имени
func (pg *PostgresRepo) GetByShortName(ctx context.Context, name string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var originalURL string

	row := pg.Conn.QueryRow(ctx,
		"SELECT original_url FROM "+tableName+" WHERE short_key=@shortKey",
		pgx.NamedArgs{"shortKey": name},
	)

	err := row.Scan(&originalURL)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return ``, err
		}

		return ``, nil
	}

	return originalURL, nil
}

// IsReady готовность репозитория
func (pg *PostgresRepo) IsReady(ctx context.Context) bool {
	err := pg.Conn.Ping(ctx)
	return err == nil
}

// RemoveByOriginalUrl удалить url
func (pg *PostgresRepo) RemoveByOriginalUrl(ctx context.Context, originalUrl string) error {

	_, err := pg.Conn.Exec(ctx,
		"DELETE FROM "+tableName+" WHERE original_url=@original_url",
		pgx.NamedArgs{"original_url": originalUrl},
	)

	return err
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) {

	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
		    id SERIAL NOT NULL PRIMARY KEY,
		    short_key varchar(`+strconv.Itoa(urlhasher.ShortNameLength)+`) UNIQUE NOT NULL,
		    original_url text NOT NULL 
		);

		CREATE UNIQUE INDEX IF NOT EXISTS `+tableName+`_original_url_unique_idx ON `+tableName+` (original_url);
	`)

	if err != nil {
		logger.Error(`create schema db error`, err)
	}
}
