package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitAccessTokens, downInitAccessTokens)
}

func upInitAccessTokens(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
			CREATE TABLE access_tokens (
				id              BIGSERIAL PRIMARY KEY,
			
				user_id         BIGINT NOT NULL,
				access_token    TEXT NOT NULL,
				refresh_token   TEXT NOT NULL,
			
				revoked_at      TIMESTAMPTZ,
			
				created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
				updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
			
				CONSTRAINT access_tokens_user_fk
					FOREIGN KEY (user_id)
					REFERENCES users(id)
					ON DELETE CASCADE
			);

			`); err != nil {
		return err
	}
	return nil
}

func downInitAccessTokens(_ context.Context, _ *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
