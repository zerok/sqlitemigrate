# sqlitemigrate

Simple migration library for SQLite. This implementation doesn't use a table to
track what version the schema is at but stores that in the `user_version`
pragma of the database itself.

## Usage

```
package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	mig "github.com/zerok/sqlitemigrate"
)

func main() {
	ctx := context.Background()

	reg := mig.NewRegistry()
	reg.RegisterMigration([]string{
		`CREATE TABLE users (id integer primary key autoincrement)`,
	}, []string{})

	db, _ := sql.Open("sqlite3", "test.sqlite")
	defer db.Close()
	if err := reg.Apply(ctx, db); err != nil {
		log.Fatal("Failed to apply migration: %s", err.Error())
	}
}
```
