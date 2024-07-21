package config

import (
	"fmt"
	"net"
	"strconv"
)

// DB contains the settings for the database connection.
type DB struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   Secret `json:"user"`
	Pass   Secret `json:"pass"`
	DBName string `json:"db_name"`
}

func (d DB) String() string {
	return fmt.Sprintf(
		"{Host: %s, Port: %d, User: %s, Pass: %s, DBName: %s}",
		d.Host,
		d.Port,
		d.User,
		d.Pass,
		d.DBName,
	)
}

// DSN returns the data source name for the database connection. It is used by the pgx library.
// It is appended with sslmode=disable option.
func (d DB) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host,
		d.Port,
		d.User,
		d.Pass,
		d.DBName,
	)
}

// ConnectionString returns a connection string to the database in a canonical form.
// It is appended with sslmode=disable option.
func (d DB) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		string(d.User),
		string(d.Pass),
		net.JoinHostPort(d.Host, strconv.Itoa(d.Port)),
		d.DBName,
	)
}
