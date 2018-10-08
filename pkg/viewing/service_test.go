package viewing_test

import (
	"testing"
	"gophoact/pkg/storage"
	"gophoact/pkg/viewing"
)

const testDbPath = "../../testdb"
const testFilepath = "../../testdata"
const testFile = "../../sampledata/TESTIMG.JPG"

func TestListAllMedia(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	a := viewing.NewService(s, fs)

	items, err := a.ListAll()
	if len(items) == 0 {
		t.Errorf("no items returned %v", err)
	} 
}