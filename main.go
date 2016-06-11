package main

import (
	"github.com/crob1140/codewiz/log"
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

	log.Debug("Opening database connection", log.Fields{"driver" : dbDriver, "dsn" : dbDSN})
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

	log.Info("Server is now listening for requests", log.Fields{"port" : port})
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
}

func assertConfigExists(key string, value string) {
	if value == "" {
		log.Fatal("Missing environment variable.", log.Fields{
			"variable": config.GetEnvironmentVariableName(key),
		});
	}
}
