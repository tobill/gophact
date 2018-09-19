package backend

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/satori/go.uuid"
)


// IMedia - general interface for mediatypes
type IMedia interface {
	Save() error
}

const testPath = "../testdb"

// MediaFile - general struct for media
type MediaFile struct {
	FileID uuid.UUID
	ID int
	Key string
	Size uint64
	Filename string
}

// LoadAll - returns all mediafiles
 func LoadAll() ([]*MediaFile, error) {
	return nil, nil
 }

// Save - save data to db
 func (mf *MediaFile) Save() error {
	dbClient, err := GetDbClient()
	if err != nil {
		return err
	}
	err = dbClient.Update(func(txn *badger.Txn) error {
		d, errint := mf.MarshalBinary()
		if errint!= nil {
			return errint
		}
		key := []byte(mf.Key)
		errint = txn.Set(key, d)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
 }

 // MarshalBinary - marshall obeject to bytearray, using json encoder
 func (mf *MediaFile) MarshalBinary() ([]byte, error) {
	return json.Marshal(mf)
 }

 // UnmarshalBinary - unmarshal from bytearray to obejct using json encoder
 func (mf *MediaFile) UnmarshalBinary(d []byte) (error) {
	err = json.Unmarshal(d, mf)
	return err
 }