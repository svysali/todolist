package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"

	// this is required to register the postgres driver.
	_ "github.com/lib/pq"
)

var conn *sqlx.DB

// GetConfig loads the config object from env vars and returns it
func GetConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("db", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// InitDb initializes the database connection from the environment variables
func InitDb() error {
	if log.GetLevel() >= log.DebugLevel {
		log.Debug("InitDb: Start")
	}

	config, err := GetConfig()
	if err != nil {
		log.WithError(err).Error("Unable to configure the database")
		return err
	}

	if conn, err = sqlx.Connect("postgres", config.ConnectionString()); err != nil {
		conn = nil
		log.WithError(err).Error("Unable to connect to the database")
		return err
	}

	if config.MaxConnections > 0 {
		conn.SetMaxOpenConns(config.MaxConnections)
	}

	if config.MaxIdleConnections > 0 {
		conn.SetMaxIdleConns(config.MaxIdleConnections)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	count, err := migrate.Exec(conn.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.WithError(err).Error("Unable to migrate the database")
		log.Fatal(err)
	}
	if count > 0 {
		log.WithField("count", count).Info("Executed migrations")
	}

	if log.GetLevel() >= log.DebugLevel {
		log.Debug("InitDb: Completed")
	}

	return nil
}

// SetDBConn is really only used to set the mockdb
func SetDBConn(db *sqlx.DB) {
	conn = db
}

// GetDBConn get the active connection
func GetDBConn() *sqlx.DB {
	if conn == nil {
		err := errors.New("Database connection not initialized. You must call `InitDB()` first")
		log.Fatal(err)
	}
	return conn
}
