-- +goose Up
ALTER TABLE short_url ADD user_id varchar(36) NOT NULL;

DROP INDEX short_url_original_url_unique_idx;
CREATE UNIQUE INDEX short_url_user_id_original_url_unique_idx ON short_url (user_id, original_url);
CREATE INDEX short_url_user_id_idx ON short_url (user_id);

-- +goose Down
DROP INDEX short_url_user_id_original_url_unique_idx;
CREATE UNIQUE INDEX short_url_original_url_unique_idx ON short_url (original_url);
ALTER TABLE short_url DROP user_id;