package speech

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	_ "github.com/go-sql-driver/mysql"
)

type Speech struct {
	dbConn *sql.DB
}

var instance *Speech

func Init() (*Speech, error) {
	if instance != nil {
		return instance, nil
	}
	db, err := connect()
	if err != nil {
		return nil, err
	}
	instance = &Speech{
		dbConn: db,
	}
	return instance, nil
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("lbry:%s@tcp(localhost)/lbry?parseTime=true", os.Getenv("MYSQL_LBRY_PASSWORD")))
	return db, errors.Err(err)
}
