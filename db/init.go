package db

import (
	"database/sql"
	"fmt"

	"voidwalker/configs"
	"voidwalker/migration"

	_ "github.com/go-sql-driver/mysql" // import mysql
	"github.com/lbryio/lbry.go/v2/extras/errors"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
)

// Init initializes a database connection based on the dsn provided. It also sets it as the global db connection.
func Init(debug bool) (*QueryLogger, error) {
	conf := configs.Configuration
	dbConn, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?parseTime=1&collation=utf8mb4_unicode_ci",
		conf.Voidwalker.User,
		conf.Voidwalker.Password,
		conf.Voidwalker.Host,
		conf.Voidwalker.Database,
	))
	if err != nil {
		return nil, errors.Err(err)
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, errors.Err(err)
	}

	logWrapper := &QueryLogger{DB: dbConn}
	if debug {
		boil.DebugMode = true
	}

	boil.SetDB(dbConn)

	migrations := &migrate.AssetMigrationSource{
		Asset:    migration.Asset,
		AssetDir: migration.AssetDir,
		Dir:      "migration",
	}
	n, migrationErr := migrate.Exec(dbConn, "mysql", migrations, migrate.Up)
	if migrationErr != nil {
		return nil, errors.Err(migrationErr)
	}
	log.Printf("Applied %d migrations", n)

	return logWrapper, nil
}

func dbInitConnection(dsn string, driverName string, debug bool) (*sql.DB, *QueryLogger, error) {
	dsn += "?parseTime=1&collation=utf8mb4_unicode_ci"
	dbConn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, errors.Err(err)
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, nil, errors.Err(err)
	}

	logWrapper := &QueryLogger{DB: dbConn}
	if debug {
		logWrapper.Logger = log.StandardLogger()
		//boil.DebugMode = true // this just prints everything twice
	}

	return dbConn, logWrapper, nil
}

// CloseDB is a wrapper function to allow error handle when it is usually deferred.
func CloseDB(db *QueryLogger) {
	if err := db.Close(); err != nil {
		log.Error("Closing DB Error: ", err)
	}
}
