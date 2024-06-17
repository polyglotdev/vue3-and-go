package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
)

// DB is a struct that represents a connection to the Postgres database
type DB struct {
	// SQL is a pointer to the database connection
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// ConnectPostgres connects to the Postgres database and returns a connection to the database
//
// Parameters:
//   - dsn: The DSN (Data Source Name) of the database to connect to.
//
// Returns:
//   - A connection to the database, or an error if the connection fails.
func ConnectPostgres(dsn string) (*DB, error) {
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(d)

	if err == nil {
		dbConn.SQL = d
	}
	return dbConn, err
}

// testDB tests the connection to the database and logs a message if the connection is successful
//
// Parameters:
//   - d: The database connection to test.
//
// Returns:
//   - An error if the connection test fails.
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		log.Errorf("Error while pinging database: %v", err)
		return err
	} else {
		log.Infof("*** Pinged database successfully! ***")
	}
	return nil
}
