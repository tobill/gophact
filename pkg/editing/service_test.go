package editing_test

import (
	"gophoact/pkg/editing"
	"gophoact/pkg/storage"
	"testing"
)

const testDbPath = "../../testdb/testbolt.db"
const testFilepath = "../../testdata"
const testFile = "../../sampledata/TESTIMG.JPG"
const testIndexPath = "../../testdbindex"

func TestLoadAndSaveID(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}
	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()

	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	if err != nil {
		t.Fatal(err)
	}
	a := editing.NewService(s, fs, is)

	var id uint64 = 1
	item, err := a.LoadMedia(id)
	if err != nil {
		t.Error(err)
	}

	if item.ID != id {
		t.Errorf("no items returned %v", err)
	}
	err = a.SaveMedia(item)
	if err != nil {
		t.Error(err)
	}

}
