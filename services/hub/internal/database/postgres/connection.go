package postgres

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

// ConnOptions represents database connection parameters.
type ConnOptions struct {
	Host           string
	Port           string
	Username       string
	DBName         string
	Password       string
	DisableSSLMode bool
}

// Connect creates a new DB connection to a PostgreSQL database.
func Connect(options ConnOptions) (*sqlx.DB, error) {
	// Convert the ConnOptions to a connection string
	connString := optionsToConnectionString(options)
	return sqlx.Connect("postgres", connString)
}

// optionsToConnectionString takes in a ConnOptions struct and returns a
// string with a formatted Postgres connection string.
func optionsToConnectionString(options ConnOptions) string {
	var sb strings.Builder

	if len(options.Host) > 1 {
		fmt.Fprintf(&sb, "host=%s ", options.Host)
	}

	if len(options.Port) > 1 {
		fmt.Fprintf(&sb, "port=%s ", options.Port)
	}

	if len(options.Username) > 1 {
		fmt.Fprintf(&sb, "user=%s ", options.Username)
	}

	if len(options.Password) > 1 {
		fmt.Fprintf(&sb, "password=%s ", options.Password)
	}

	if len(options.DBName) > 1 {
		fmt.Fprintf(&sb, "dbname=%s ", options.DBName)
	}

	if options.DisableSSLMode {
		fmt.Fprint(&sb, "sslmode=disable ")
	}

	return sb.String()
}
