package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for a db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
		CREATE TABLE users (
			user_id UUID,
			name TEXT,
			email TEXT UNIQUE,
			password_hash TEXT,

			date_created TIMESTAMP,
			date_updated TIMESTAMP,

			PRIMARY KEY (user_id)
		);`,
	},
}
