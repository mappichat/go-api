package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Sqldb *sqlx.DB

func SqlInitialize(connectString string) (*sqlx.DB, error) {
	var err error
	Sqldb, err = sqlx.Connect("postgres", connectString)
	if err != nil {
		return Sqldb, err
	}
	if err = Sqldb.Ping(); err != nil {
		return Sqldb, err
	}
	return Sqldb, nil
}
