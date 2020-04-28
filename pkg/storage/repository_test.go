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


const testDbPath = "../../testdb/1q"
const testFilepath = "../../testdata"
const testFile = "../../sampledata/TESTIMG.JPG"

func TestMediaAdd(t *testing.T){
	r, err := storage.NewDbStorage(testDbPath)
	defer r.CloseDb()
	if err != nil { t.Error(err) }
	m := adding.Media{
		Filename: "Tddddestfile",
		Size: 12345,
	}
	m.FileID = uuid.NewV4()
	_, err= r.AddMedia(&m)
	if err != nil { t.Error(err) }


}

func TestFileAdd(t *testing.T) {
	r := storage.NewFileStorage(testFilepath)
	fp, err := os.Open(testFile)
	defer fp.Close()
	if err != nil {
		t.Error(err)
	}
	m := adding.Media{
		Filename: "Tesddddtfile",
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
	defer db.CloseDb()
	if err != nil { t.Error(err) }

	items, err := db.ListAll(0,30)
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
	defer db.CloseDb()

	if err != nil { t.Error(err) }
	
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

	fileID := "05e8b739-96c5-4376-8a8e-cbfd54c74348"
	//var r *os.File
	r, err := f.GetOriginalFile(fileID)
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

	r, err := f.GetOriginalFile(fileID)
	defer r.Close()

	if err == nil {
		t.Errorf("error should be returned")
	}

	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("error should be os.PathError")
	}
}

func TestSaveMedia(t *testing.T) {
	db, err := storage.NewDbStorage(testDbPath)
	defer db.CloseDb()

	if err != nil { t.Error(err) }

	var objID uint64 = 1

	m, err := db.LoadMedia(objID)
	if err != nil { t.Error(err) }

	log.Printf("%v", m)

	err = db.SaveMedia(m)
	if err != nil { t.Error(err) }

	mfs, err := db.GetMediaPerMimetype(m.MimeType.MIME.Value)

	for x := range mfs {
		log.Printf("%v",mfs[x])
	}

	if err != nil {
		t.Error(err)
	}
}



