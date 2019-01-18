package adding_test

import (
	"gophoact/pkg/jobqueue"
	"gophoact/pkg/editing"
	"io"
	"os"
	"testing"
	"gophoact/pkg/storage"
	"gophoact/pkg/adding"
	"mime/multipart"
)

const testDbPath = "../../testdb/testbolt.db"
const testFilepath = "../../testdata"
const testIndexPath = "../../testdbindex"
const testFile = "../../sampledata/TESTIMG.JPG"

func TestAddMedia(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	is, err := storage.NewIndexStorage(testIndexPath)
	if err != nil { t.Error(err) }
	defer is.CloseIndex()

	e := editing.NewService(s, fs, is)

	jq := jobqueue.NewService(e)
	defer jq.CloseQueue()

	a := adding.NewService(s, fs, jq)

	var size int64
	size = 543455

	
    filename := "testfile.jpg"
	mph := multipart.FileHeader{
		Filename: filename, 
		Size: size,
	}

	filepath := "../../sampledata/TESTIMG.JPG"
	mpf, err := os.Open(filepath)
	defer mpf.Close()

	var mpr io.Reader = mpf
	
	if err != nil { t.Error(err) }
	_, err = a.AddMedia(&mpr, &mph)

	if err != nil { t.Error(err) }

}