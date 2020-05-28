# sqlitemigrate

<a href="https://pkg.go.dev/github.com/zerok/sqlitemigrate?tab=doc"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&amp;logoColor=white&amp;style=flat-square" alt="go.dev reference"></a>

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
