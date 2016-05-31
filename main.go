package main

import (
	"log"
	"database/sql"
	"github.com/crob1140/codewiz/config"
	"github.com/crob1140/codewiz/datastore"
	_ "github.com/mattn/go-sqlite3"
)
func main() {

	dbDriver, err := config.GetString("database.driver")
	if err != nil {
		log.Fatal(err)
	}

	dbHost, err := config.GetString("database.host")
	if err != nil {
		log.Fatal(err)
	}

	dbPort, err := config.GetString("database.port")
	if err != nil {
		log.Fatal(err)
	}

	host, err := config.GetString("host")
	if err != nil {
		log.Fatal(err)
	}

	port, err := config.GetString("port")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(dbDriver, dbHost + ":" + dbPort)
	if err != nil {
		log.Fatal(err)
	}

	ds := datastore.NewDatastore(db , dbDriver)
	server := NewServer(ds)
	server.ListenAndServe(host + ":" + port)
}

func getDialectForDriverName() {

}
