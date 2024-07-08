package config

import (
	"fmt"
	"net"
	"strconv"
)

// Cache contains the settings for the connection to application cache.
type Cache struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	User   Secret `json:"user"`
	Pass   Secret `json:"pass"`
	DBName int    `json:"db_name"`
}

func (d Cache) String() string {
	return fmt.Sprintf(
		"{Host: %s, Port: %d, User: %s, Pass: %s, DBName: %d}",
		d.Host,
		d.Port,
		d.User,
		d.Pass,
		d.DBName,
	)
}

func (d Cache) ConnectionString() string {
	return fmt.Sprintf(
		"redis://%s:%s@%s/%d",
		string(d.User),
		string(d.Pass),
		net.JoinHostPort(d.Host, strconv.Itoa(d.Port)),
		d.DBName,
	)
}
