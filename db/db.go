package db

import (
	"context"
	"database/sql"
	"time"

	pq "github.com/lib/pq"
	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Database is an interface for database
type Database interface {
	Connect(connStr string) (*sql.DB, error)
	RunMigration(connStr string) error
	RollbackMigration(connStr string) error
}

// RealDatabase is a real implementation of Database
type RealDatabase struct {
	migrateDir string
}

// NewDatabase returns a new instance of Database
func NewDatabase(migrateDir string) Database {
	return &RealDatabase{migrateDir}
}

// Connect to database using connection string
func (rdb *RealDatabase) Connect(connStr string) (*sql.DB, error) {
	const maxOpenConns = 100
	const maxIdleConns = 10

	connector, err := pq.NewConnector(connStr)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Minute * 1)

	// Check if database is alive
	if err = db.PingContext(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}

// RunMigration runs migration
func (rdb *RealDatabase) RunMigration(connStr string) error {
	db, err := rdb.Connect(connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	migration, err := migrate.New(rdb.migrateDir, connStr)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil {

		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	defer migration.Close()
	return nil
}

func (rdb *RealDatabase) RollbackMigration(connStr string) error {
	migration, err := migrate.New(rdb.migrateDir, connStr)
	if err != nil {
		return err
	}

	if err := migration.Down(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}

	defer migration.Close()
	return nil
}
