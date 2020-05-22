package deleting_test

import (
	"gophoact/pkg/deleting"
	"gophoact/pkg/storage"
	"testing"
)

const testDbPath = "../../testdb/testbolt.db"
const testFilepath = "../../testdata"
const testIndexPath = "../../testdbindex"
const testFile = "../../sampledata/TESTIMG.JPG"


func TestService_DeleteMedia(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath)
	if err != nil { t.Fatal(err)	}
	is, err := storage.NewIndexStorage(testIndexPath)
	if err != nil { t.Error(err) }
	defer is.CloseIndex()

	deleting.NewService(s, fs, is)
}
