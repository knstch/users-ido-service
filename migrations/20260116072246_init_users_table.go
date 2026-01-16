package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitUsersTable, downInitUsersTable)
}

func upInitUsersTable(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
				CREATE TABLE users (
					id              BIGSERIAL PRIMARY KEY,
				
					google_sub      TEXT NOT NULL,
					email           TEXT NOT NULL,
					first_name       TEXT,
					last_name        TEXT,
					profile_picture  TEXT,
				
					created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
					updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
				    deleted_at 		TIMESTAMPTZ NULL
				);
				
				CREATE UNIQUE INDEX users_google_sub_uk
					ON users (google_sub);
				
				CREATE UNIQUE INDEX users_email_uk
					ON users (email);
`); err != nil {
		return err
	}
	return nil
}

func downInitUsersTable(_ context.Context, _ *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
