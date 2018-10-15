package storage_test

import (
	"log"
	"os"
	"io"
	"github.com/satori/go.uuid"
	"testing"
	"gophoact/pkg/storage"
	"gophoact/pkg/adding"
)


const testDbPath = "../../testdb"
const testFilepath = "../../testdata"
const testFile = "../../sampledata/TESTIMG.JPG"

func TestMediaAdd(t *testing.T){
	r, err := storage.NewDbStorage(testDbPath)
	defer r.CloseDb()
	if err != nil { t.Error(err) }
	m := adding.Media{
		Filename: "Testfile",
		Size: 12345,
	}
	m.FileID = uuid.NewV4()
	r.AddMedia(&m)
}

func TestFileAdd(t *testing.T) {
	r := storage.NewFileStorage(testFilepath)
	fp, err := os.Open(testFile)
	defer fp.Close()
	if err != nil {
		t.Error(err)
	}
	m := adding.Media{
		Filename: "Testfile",
		Size: 12345,
		FileID: uuid.NewV4(),
	}
	var fpr io.Reader = fp
	err = r.AddFile(&fpr, &m)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAllMedia(t *testing.T) {
	
	db, err := storage.NewDbStorage(testDbPath)
	if err != nil { t.Error(err) }
	defer db.CloseDb()

	items, err := db.ListAll()
	if err != nil { t.Error(err) }

	if len(items) == 0 {
		t.Errorf("no items returned")
	}

	for _, ele := range items {
		log.Printf("%v", ele)
	
	}
}

func TestGetById(t *testing.T) {
	
	db, err := storage.NewDbStorage(testDbPath)
	if err != nil { t.Error(err) }
	defer db.CloseDb()
	
	var id uint64 = 1 

	item, err := db.GetByID(id)
	if err != nil { 
		t.Error(err) 
	}

	if item.ID != 1 {
		t.Errorf("item not found")
	}
	log.Printf("%v", item)

}

func TestGetFileByFileId(t *testing.T) {
	
	f := storage.NewFileStorage(testFilepath)

	fileID := "a1b3f162-13df-44c2-a064-116537443b80"
	//var r *os.File
	r, err := f.GetFile(fileID)
	defer r.Close()
	if err != nil { t.Error(err) }

	b := make([]byte, 5)
	n1, err := r.Read(b)
	if err != nil { t.Error(err) }
	log.Printf("%v", n1)

}

func TestGetFileByFileIdNotFound(t *testing.T) {
	f := storage.NewFileStorage(testFilepath)

	fileID := "000"
	//var r *os.File

	r, err := f.GetFile(fileID)
	defer r.Close()

	if err == nil {
		t.Errorf("error should be returned")
	}

	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("error should be os.PathError")
	}


}


