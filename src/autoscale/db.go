package autoscale

import (
	"database/sql"
	"net/url"

	// import db driver
	_ "github.com/lib/pq"
)

// NewDB creates a db connection.
func NewDB(user, password, addr, database string) (*sql.DB, error) {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   addr,
		Path:   database,
	}

	v := url.Values{}
	v.Set("sslmode", "disable")

	u.RawQuery = v.Encode()

	return sql.Open("postgres", u.String())
}
