package repository

import (
	"context"
	"strconv"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pgRepo *Postgres

const tableName = `short_url`

// Postgres connection structure
type Postgres struct {
	Conn *pgxpool.Pool
}

// Add Добавить URL
func (pg *Postgres) Add(ctx context.Context, name string) (string, error) {
	return ``, nil
}

// GetByShortName получить URL по короткому имени
func (pg *Postgres) GetByShortName(ctx context.Context, name string) (string, error) {
	return ``, nil
}

// IsReady готовность репозитория
func (pg *Postgres) IsReady(ctx context.Context) bool {
	err := pg.Conn.Ping(ctx)
	return err == nil
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) {

	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
		    id SERIAL NOT NULL PRIMARY KEY,
		    user_id varchar(36) NOT NULL,
		    short_key varchar(`+strconv.Itoa(urlhasher.ShortNameLength)+`) NOT NULL,
		    original_url text NOT NULL,
			is_deleted boolean NOT NULL DEFAULT false
		);

		CREATE UNIQUE INDEX IF NOT EXISTS `+tableName+`_user_id_short_key_is_deleted_unique_idx ON `+tableName+` (user_id, short_key, is_deleted);
		CREATE UNIQUE INDEX IF NOT EXISTS `+tableName+`_user_id_original_url_is_deleted_unique_idx ON `+tableName+` (user_id, original_url, is_deleted);
	`)

	if err != nil {
		logger.Error(`create schema db error`, err)
	}
}
