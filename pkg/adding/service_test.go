package adding_test

import (
	"gophoact/pkg/adding"
	"gophoact/pkg/storage"
	"gophoact/pkg/viewing"
	"io"
	"log"
	"mime/multipart"
	"os"
	"testing"
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

	//docs, err := is.FindDocuments("39a641f0d33003a5a729a8c2415f4b0b130e4620")

	a := adding.NewService(s, fs, is)

	vs := viewing.NewService(s, fs)

	d, err := vs.ListAll(0,20)
	for x:=0; x< len(d); x++ {
		log.Printf("%v", d[x])
	}


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

	if err != nil {
		t.Error(err)
	}

	if _, ok := err.(*adding.DuplicateFileError); !ok {
		t.Errorf("file should already be there")
	}

	//if err != nil { t.Error(err) }

}
