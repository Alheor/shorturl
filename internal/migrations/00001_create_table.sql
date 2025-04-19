-- +goose Up
CREATE TABLE IF NOT EXISTS short_url (
	id SERIAL NOT NULL PRIMARY KEY,
	short_key varchar(8) UNIQUE NOT NULL,
	original_url text NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS short_url_original_url_unique_idx ON short_url (original_url);

-- +goose Down
DROP INDEX short_url_original_url_unique_idx;
DROP TABLE short_url;
