package storage_test

import (
	"testing"
	"gophoact/pkg/storage"

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
	err = is.AddDocument(item)
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

