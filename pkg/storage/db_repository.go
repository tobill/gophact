package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"gophoact/pkg/adding"
	"gophoact/pkg/editing"
	"gophoact/pkg/viewing"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

//DbStorage stores data in database
type DbStorage struct {
	dbClient *bolt.DB
}

const mediaKeyPrefix = "media:"

var mediaBucket = []byte("media")
var mimetypeBucket = []byte("mimetype")
var mediaKeyCounter = []byte("mediaKeyCounter")

const mimetypeKeyPrefix = "mimetype:"

//NewDbStorage create new storage for data
func NewDbStorage(dbPath string) (*DbStorage, error) {
	dbClient, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
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

func getBucket(bucketname []byte, tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket(bucketname)
	if bucket != nil {
		return bucket, nil
	}

	bucket, err := tx.CreateBucket(bucketname)
	if err != nil {
		return nil, err
	}
	return bucket, err
}

//AddMedia inserts data into db
func (s *DbStorage) AddMedia(media *adding.Media) (uint64, error) {

	err := s.dbClient.Update(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mediaBucket, txn)
		if err != nil {
			return err
		}
		if media.Key == "" {
			id, err := bucket.NextSequence()
			if err != nil {
				return err
			}
			media.ID = id
			media.Key = mediaKeyPrefix + strconv.FormatUint(media.ID, 10)
		}
		// convert to storage model
		sMedia := &Media{
			FileID:   media.FileID,
			Filename: media.Filename,
			Size:     media.Size,
			ID:       media.ID,
			Key:      media.Key,
		}

		d, errint := sMedia.marshalMedia()
		if errint != nil {
			return errint
		}
		log.Printf("add data")
		key := []byte(media.Key)
		errint = bucket.Put(key, d)
		log.Printf("%v", errint)
		return errint
	})
	if err != nil {
		return 0, err
	}
	return media.ID, nil
}

//ListAll returns all entries
func (s *DbStorage) ListAll(numStart uint64, numItems uint64) ([]*viewing.Media, error) {
	var mfs []*viewing.Media
	log.Printf("starting")
	err := s.dbClient.View(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mediaBucket, txn)
		if err != nil {
			return err
		}
		it := bucket.Cursor()
		bucket.Stats()
		for k, v := it.First(); k != nil; k, v = it.Next() {
			sMedia := Media{}
			sMedia.unmarshalMedia(v)
			// convert to storage model
			vMedia := &viewing.Media{
				FileID:   sMedia.FileID,
				Filename: sMedia.Filename,
				Size:     sMedia.Size,
				ID:       sMedia.ID,
				Key:      sMedia.Key,
				Mimetype: sMedia.Mimetype,
			}
			mfs = append(mfs, vMedia)
		}
		return nil
	})
	return mfs[numStart : numStart+numItems], err
}

func (m *Media) marshalMedia() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(m)
	return b.Bytes(), err
}

func (m *Media) unmarshalMedia(d []byte) error {
	b := bytes.NewBuffer(d)
	dec := gob.NewDecoder(b)
	err := dec.Decode(m)
	return err
}

//GetByID returns  specific data by id
func (s *DbStorage) GetByID(id uint64) (*viewing.Media, error) {
	var media viewing.Media
	err := s.dbClient.View(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mediaBucket, txn)
		if err != nil {
			return err
		}

		sMedia := Media{}
		key := mediaKeyPrefix + strconv.FormatUint(id, 10)
		item := bucket.Get([]byte(key))
		sMedia.unmarshalMedia(item)
		media = viewing.Media{
			FileID:   sMedia.FileID,
			Filename: sMedia.Filename,
			Size:     sMedia.Size,
			ID:       sMedia.ID,
			Key:      sMedia.Key,
			Mimetype: sMedia.Mimetype,
		}
		return nil
	})
	return &media, err
}

//LoadMedia returns  specific data by id
func (s *DbStorage) LoadMedia(id uint64) (*editing.Media, error) {
	var media editing.Media
	err := s.dbClient.View(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mediaBucket, txn)
		if err != nil {
			return err
		}

		sMedia := Media{}
		key := mediaKeyPrefix + strconv.FormatUint(id, 10)
		item := bucket.Get([]byte(key))
		if item == nil {
			return errors.New("key not found")
		}
		sMedia.unmarshalMedia(item)
		media = editing.Media{
			FileID:   sMedia.FileID,
			Filename: sMedia.Filename,
			Size:     sMedia.Size,
			ID:       sMedia.ID,
			Key:      sMedia.Key,
			MimeType: sMedia.Mimetype,
			Versions: sMedia.Versions,
		}
		return nil
	})
	return &media, err
}

//SaveMedia saves mediadata
func (s *DbStorage) SaveMedia(media *editing.Media) error {
	// convert to storage model
	sMedia := &Media{
		FileID:   media.FileID,
		Filename: media.Filename,
		Size:     media.Size,
		ID:       media.ID,
		Key:      media.Key,
		Mimetype: media.MimeType,
		Versions: media.Versions,
	}

	err := s.dbClient.Update(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mediaBucket, txn)
		if err != nil {
			return err
		}

		d, errint := sMedia.marshalMedia()
		if errint != nil {
			return errint
		}
		key := []byte(media.Key)
		errint = bucket.Put(key, d)
		return errint
	})
	if err != nil {
		return err
	}
	err = s.updateMimetypeIndex(media.MimeType.MIME.Value, media.ID)
	if err != nil {
		return err
	}
	return nil
}

func findValPos(slice []uint64, val uint64) int {
	for x := range slice {
		if slice[x] > val {
			return x
		}
		if slice[x] == val {
			return -1
		}

	}
	return len(slice)
}

func insertIntoList(slice []uint64, val uint64) []uint64 {
	length := len(slice)
	if length == 0 {
		return append(slice, val)
	}
	pos := findValPos(slice, val)
	if pos == -1 {
		return slice
	}
	slice = append(slice, 1)
	for i := length; i > pos; i-- {
		slice[i] = slice[i-1]
	}
	slice[pos] = val
	return slice
}

func (s *DbStorage) updateMimetypeIndex(mt string, objID uint64) error {
	err := s.dbClient.Update(func(txn *bolt.Tx) error {
		bucket, err := getBucket(mimetypeBucket, txn)
		if err != nil {
			return err
		}

		key := []byte(mimetypeKeyPrefix + mt)
		item := bucket.Get(key)
		var mti MimetypeIndex
		if item == nil {
			mti = MimetypeIndex{
				ObjectIds: append(make([]uint64, 0), objID),
			}
		} else {
			errint := mti.unmarshalMimetypeIndex(item)
			if errint != nil {
				return errint
			}
			mti.ObjectIds = insertIntoList(mti.ObjectIds, objID)

		}

		b, errint := mti.marshalMimetypeIndex()
		if errint != nil {
			return errint
		}
		errint = bucket.Put(key, b)
		return errint
	})
	return err
}

//GetMediaPerMimetype returns all media per mimetype
func (s *DbStorage) GetMediaPerMimetype(mt string) ([]*viewing.Media, error) {
	var mfs []*viewing.Media
	err := s.dbClient.View(func(txn *bolt.Tx) error {
		objectIds := MimetypeIndex{}
		key := []byte(mimetypeKeyPrefix + mt)
		bucket, err := getBucket(mimetypeBucket, txn)
		if err != nil {
			return err
		}
		item := bucket.Get(key)
		errint := objectIds.unmarshalMimetypeIndex(item)
		if errint != nil {
			return errint
		}
		for i := range objectIds.ObjectIds {
			media, errint := s.GetByID(objectIds.ObjectIds[i])
			if errint != nil {
				return errint
			}
			mfs = append(mfs, media)
		}
		return errint
	})

	return mfs, err
}

func (mt *MimetypeIndex) marshalMimetypeIndex() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mt)
	return b.Bytes(), err
}

func (mt *MimetypeIndex) unmarshalMimetypeIndex(d []byte) error {
	b := bytes.NewBuffer(d)
	dec := gob.NewDecoder(b)
	err := dec.Decode(mt)
	return err
}
