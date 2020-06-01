package sqlitemigrate

import (
	"context"
	"database/sql"
	"fmt"
)

// MigrationRegistry is the place where all migrations attached.
type MigrationRegistry struct {
	migrations []Migration
}

// NewRegistry creates a new MigrationRegistry.
func NewRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrations: make([]Migration, 0, 10),
	}
}

// GetCurrentVersion returns the version of the database.
func (r *MigrationRegistry) GetCurrentVersion(ctx context.Context, db *sql.DB) (int, error) {
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return -1, err
	}
	return version, nil
}

// Apply tries to run all unapplied migrations onto the database.
func (r *MigrationRegistry) Apply(ctx context.Context, db *sql.DB) error {
	current, err := r.GetCurrentVersion(ctx, db)
	if err != nil {
		return err
	}
	for idx, mig := range r.migrations {
		mig.Version = idx + 1
		if current >= idx+1 {
			continue
		}
		if err := mig.Apply(ctx, db); err != nil {
			return err
		}
	}
	return nil
}

// Reverse tries to reverse the last applied migration.
func (r *MigrationRegistry) Reverse(ctx context.Context, db *sql.DB) error {
	current, err := r.GetCurrentVersion(ctx, db)
	if err != nil {
		return err
	}
	if current == 0 {
		return nil
	}
	if current > len(r.migrations) {
		return fmt.Errorf("current version number does not match known migration")
	}
	r.migrations[current-1].Version = current
	return r.migrations[current-1].Reverse(ctx, db)
}

// RegisterMigration adds a new set of up- and down-statements to the registry.
func (r *MigrationRegistry) RegisterMigration(ups []string, downs []string) {
	m := Migration{
		Up:   ups,
		Down: downs,
	}
	r.migrations = append(r.migrations, m)
}

// Migration is basically just a list of SQL statements that should be run to
// upgrade or to downgrade the database schema.
type Migration struct {
	Version int
	Up      []string
	Down    []string
}

// Apply runs the migration's Up statements onto the given DB in a dedicated
// transaction.
func (m *Migration) Apply(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, up := range m.Up {
		if _, err := tx.ExecContext(ctx, up); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", m.Version)); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (m *Migration) Reverse(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, down := range m.Down {
		if _, err := tx.ExecContext(ctx, down); err != nil {
			tx.Rollback()
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", m.Version-1)); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
