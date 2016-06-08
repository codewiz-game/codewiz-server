package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/crob1140/codewiz/config"
	"github.com/crob1140/codewiz/config/keys"
	"github.com/crob1140/codewiz/datastore"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	initLogger()

	dbDriver := config.GetString(keys.DatabaseDriver)
	assertConfigExists(keys.DatabaseDriver, dbDriver)

	dbDSN := config.GetString(keys.DatabaseDSN)
	assertConfigExists(keys.DatabaseDSN, dbDSN)

	port := config.GetString(keys.Port)
	assertConfigExists(keys.Port, port)

	ds, err := datastore.Open(dbDriver, dbDSN)
	if err != nil {
		log.Fatal(err)
	}

	errs, ok := ds.UpSync()
	if !ok {
		for _, err := range errs {
			log.Error(err)
		}
	}

	server := NewServer(ds)
	server.ListenAndServe(":" + port)
}

func initLogger() {
	logLevel := config.GetString(keys.LogLevel)
	level := log.InfoLevel
	if logLevel != "" {
		var err error
		if level, err = log.ParseLevel(logLevel); err != nil {
			log.Fatal(err)
		}
	}

	log.SetLevel(level)
	log.SetFormatter(&log.JSONFormatter{})
}

func assertConfigExists(key string, value string) {
	if value == "" {
		log.WithFields(log.Fields{
			"variable": config.GetEnvironmentVariableName(key),
		}).Fatal("Missing environment variable.")
	}
}
