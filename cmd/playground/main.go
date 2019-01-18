package main 


import (
	"log"
    "github.com/boltdb/bolt"
)


const path = "./playdb/bold.db"

//DbStorage stores data in database 
type DbStorage struct {
	dbClient *bolt.DB
}

//NewDbStorage create new storage for data
func NewDbStorage(dbPath string) (*DbStorage, error) {
	dbClient, err :=  bolt.Open(path, 0600, nil)
	if err != nil { return nil, err }
	s := DbStorage{
		dbClient: dbClient,
	}
	return &s, err
}

func openAndWrite() error {
    db, err := NewDbStorage(path)
    if err != nil {
        return err
    }
    buckentame := []byte("world")
    key := []byte("testkey")
    d := []byte("data")

    err = db.dbClient.View(func(txn *bolt.Tx) error {
        bucket := txn.Bucket(buckentame)
        if bucket == nil {
            log.Printf("Bucket %q not found!", buckentame)
            return nil
        }

        val := bucket.Get(key)      
        log.Printf("%v", val)  
        return nil
    })
	if err != nil {
		return err
    }

    err = db.dbClient.Update(func(txn *bolt.Tx) error {
        b ,errint := txn.CreateBucketIfNotExists(buckentame)
        if errint != nil {
            return errint
        }
        errint = b.Put(key, d)
		return errint
	})
	if err != nil {
		return err
    }


    db.dbClient.Close()
	return nil
}

func main() {
    err := openAndWrite()
    log.Print(err)
}