package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepo connection structure
type PostgresRepo struct {
	Conn *pgxpool.Pool
}

// Add Добавить URL
func (pg *PostgresRepo) Add(ctx context.Context, user *models.User, name string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	hash := urlhasher.GetHash(name)

	_, err := pg.Conn.Exec(ctx,
		"INSERT INTO short_url (user_id, short_key, original_url) VALUES (@userId, @shortKey, @originalURL)",
		pgx.NamedArgs{"userId": user.ID, "shortKey": hash, "originalURL": name},
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := pg.Conn.QueryRow(ctx,
				"SELECT short_key FROM short_url WHERE user_id=@userId AND original_url=@originalUrl",
				pgx.NamedArgs{"userId": user.ID, "originalUrl": name},
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
func (pg *PostgresRepo) AddBatch(ctx context.Context, user *models.User, list *[]models.BatchEl) error {

	tx, err := pg.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	var entries [][]any
	for _, v := range *list {
		entries = append(entries, []any{user.ID, v.ShortURL, v.OriginalURL})
	}

	_, err = pg.Conn.CopyFrom(
		ctx,
		pgx.Identifier{`short_url`},
		[]string{"user_id", "short_key", "original_url"},
		pgx.CopyFromRows(entries),
	)

	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

// GetByShortName получить URL по короткому имени
func (pg *PostgresRepo) GetByShortName(ctx context.Context, user *models.User, name string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var originalURL string

	row := pg.Conn.QueryRow(ctx,
		"SELECT original_url FROM short_url WHERE user_id=@userId AND short_key=@shortKey",
		pgx.NamedArgs{"userId": user.ID, "shortKey": name},
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

// RemoveByOriginalURL удалить url
func (pg *PostgresRepo) RemoveByOriginalURL(ctx context.Context, user *models.User, originalURL string) error {

	_, err := pg.Conn.Exec(ctx,
		"DELETE FROM short_url WHERE user_id=@userId AND original_url=@original_url",
		pgx.NamedArgs{"userId": user.ID, "original_url": originalURL},
	)

	return err
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) error {

	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS short_url (
		    id SERIAL NOT NULL PRIMARY KEY,
		    user_id varchar(36) NOT NULL,
		    short_key varchar(`+strconv.Itoa(urlhasher.HashLength)+`) UNIQUE NOT NULL,
		    original_url text NOT NULL 
		);

		CREATE UNIQUE INDEX IF NOT EXISTS short_url_user_id_original_url_unique_idx ON short_url (user_id, original_url);
	`)

	if err != nil {
		return err
	}

	return nil
}
