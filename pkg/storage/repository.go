package storage

import (
	"io"
	"path/filepath"
	"encoding/gob"
	"bytes"
	"strconv"
	"github.com/dgraph-io/badger"
	"os"
	"gophoact/pkg/adding"
	"gophoact/pkg/viewing"
)

//DbStorage stores data in database 
type DbStorage struct {
	dbClient *badger.DB
}

//FileStorage stores filedata on disk 
type FileStorage struct {
	dirpathOriginal string
}

const mediaKeyPrefix = "media:"
var mediaKeyCounter = []byte("mediaKeyCounter")

//NewDbStorage create new storage for data
func NewDbStorage(dbPath string) (*DbStorage, error) {
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath
	dbClient, err :=  badger.Open(opts)
	if err != nil { return nil, err }
	s := DbStorage{
		dbClient: dbClient,
	}
	return &s, err
}

// CloseDb closes link to db
func (s *DbStorage) CloseDb() error {
	err := s.dbClient.Close()
	return err
}

//NewFileStorage create new storage for files
func NewFileStorage(dirpath string) (*FileStorage) {
	s := FileStorage{
		dirpathOriginal: filepath.Join(dirpath, "original"),
	}
	return &s
}

//AddMedia inserts data into db
func (s *DbStorage) AddMedia(media *adding.Media) error {
	if media.Key == "" {
		seq, err := s.dbClient.GetSequence(mediaKeyCounter, 1)
		defer seq.Release()
		if err != nil { return err }
		media.ID, err = seq.Next()
		if err != nil { return err }
		media.Key = mediaKeyPrefix + strconv.FormatUint(media.ID, 10)
	}

	// convert to storage model
	sMedia := &Media{
		FileID: media.FileID,
		Filename: media.Filename,
		Size: media.Size,
		ID: media.ID,
		Key: media.Key,
	}

	err := s.dbClient.Update(func(txn *badger.Txn) error {
		d, errint := sMedia.marshalMedia()
		if errint!= nil {
			return errint
		}
		key := []byte(media.Key)
		errint = txn.Set(key, d)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}


//ListAll returns all entries
func (s *DbStorage) ListAll() ([]*viewing.Media, error) {
	var mfs  []*viewing.Media
	err := s.dbClient.View(func (txn *badger.Txn) error{
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		it.Close()
		for it.Rewind(); it.ValidForPrefix([]byte(mediaKeyPrefix)); it.Next() {
			sMedia := Media{}
 			item := it.Item()
			val, errint := item.Value()
			if errint != nil { return errint }
			sMedia.unmarshalMedia(val)
				// convert to storage model
			vMedia := &viewing.Media{
				FileID: sMedia.FileID,
				Filename: sMedia.Filename,
				Size: sMedia.Size,
				ID: sMedia.ID,
				Key: sMedia.Key,
			}

			mfs = append(mfs, vMedia)
		}
		return nil
	})
	return mfs, err

}

func (m *Media) marshalMedia() ([]byte, error) {
	var  b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m)
	return b.Bytes(), err
}

func (m *Media) unmarshalMedia(d []byte) (error) {
	b := bytes.NewBuffer(d)
	dec := gob.NewDecoder(b)
	err := dec.Decode(m) 
	return err
}

//AddFile inserts data into db
func (s *FileStorage) AddFile(source *io.Reader, media *adding.Media) error {
	fpath := filepath.Join(s.dirpathOriginal, media.FileID.String())
	fd, err := os.Create(fpath)
	if err != nil { return err}
	defer fd.Close()
	_, err = io.Copy(fd, *source)
	return err
}
