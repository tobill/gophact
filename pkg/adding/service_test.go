package adding_test

import (
	"io"
	"os"
	"testing"
	"gophoact/pkg/storage"
	"gophoact/pkg/adding"
	"mime/multipart"
)

const testDbPath = "../../testdb"
const testFilepath = "../../testdata"
const testFile = "../../sampledata/TESTIMG.JPG"

func TestAddMedia(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	a := adding.NewService(s, fs)

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
	a.AddMedia(&mpr, &mph)


}