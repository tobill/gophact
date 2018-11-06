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


func TestGetByID(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	a := viewing.NewService(s, fs)

	var id uint64 = 1
	item, err := a.GetByID(id)
	if err != nil {
		t.Error(err)
	} 

	if item.ID  != id {
		t.Errorf("no items returned %v", err)
	} 
}



func TestGetFileById(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	a := viewing.NewService(s, fs)

	var fid uint64 = 1

	r, err := a.GetFileByID(fid)

	if err != nil { t.Error(err) }
	defer r.Close()
	b := make([]byte, 5)
    _, err = r.Read(b)
	if err != nil { t.Error(err) }
}

func TestGetFileByIdNotFound(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	a := viewing.NewService(s, fs)

	var fid uint64 = 20

	r, err := a.GetFileByID(fid)
	defer r.Close()

	if err == nil {
		t.Errorf("error should be returned")
	}

	if _, ok :=err.(*viewing.NotFoundError); !ok  {
		t.Errorf("not found error should be returned")
	}


}

