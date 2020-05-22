package storage_test

import (
	uuid "github.com/satori/go.uuid"
	"gophoact/pkg/adding"
	"gophoact/pkg/deleting"
	"gophoact/pkg/storage"
	"testing"
)
const testIndexPath = "../../testdbindex"

func TestAddAndFindDoc(t *testing.T) {
	
	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()

	if err != nil { t.Error(err) }
	
	db, err := storage.NewDbStorage(testDbPath)
	defer db.CloseDb()
	if err != nil { t.Error(err) }
	
	var id uint64 = 1 

	item, err := db.LoadMedia(id)
	if err != nil { 
		t.Error(err) 
	}
	err = is.UpdateDocument(item)
	result, err := is.FindDocuments(item.Filename) 
	if err != nil {
		t.Error(err)
	}

	found := false
	for r := range result {
		if result[r] == item.Key {
			found = true
		}
	}
	if !found {
		t.Error("doc not found")
	}
}

func TestIndexStorage_DeleteMedia(t *testing.T) {

	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()

	if err != nil {
		t.Error(err)
	}
	const searchString =  "de63dfsf"
	m := adding.Media{
		FileID:  uuid.NewV4(),
		ID:       0265635,
		Key:      "media:0265635",
		Size:     550,
		Filename: "testfilename",
		CheckSum: searchString,
	}

	is.AddDocument(m)

	docs, err := is.FindDocuments(searchString)

	if err != nil {
		t.Error(err)
	}

	if len(docs) == 0 {
		t.Error("should find something to delete")
	}
	d1 := docs[0]

	dmedia := deleting.Media{
		 FileID: uuid.UUID{},
		 Key:    d1,
	}

	is.DeleteMedia(&dmedia)

	docs, err = is.FindDocuments(searchString)
	if err != nil {
		t.Error(err)
	}

	if len(docs) > 0 {
		t.Error("document should be deleted")
	}

}