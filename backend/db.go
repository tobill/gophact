package backend


import (
	"log"
	"github.com/dgraph-io/badger"
	"errors"
)

// DbClient holding global db
var dbClient *badger.DB
var err error
// path to db

// GetDbClient open db if needed and return instance of badger
func GetDbClient() (*badger.DB, error) {
	if dbClient == nil {
		return nil, errors.New("connection must be initialized")
	}
	return dbClient, nil
}

// Connect to client
func Connect(dbPath string) {
	if dbClient == nil {
		opts := badger.DefaultOptions
		opts.Dir = dbPath
		opts.ValueDir = dbPath
		dbClient, err =  badger.Open(opts)
		if err != nil {
			log.Fatal("error initializing db")
		}
	}
}

// Close - closes db
func Close() {
	dbClient.Close()
}