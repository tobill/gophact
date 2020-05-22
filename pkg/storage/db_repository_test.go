package storage_test

import (
	"gophoact/pkg/deleting"
	"gophoact/pkg/storage"
	"testing"
)

import (
	"github.com/satori/go.uuid"
	"gophoact/pkg/adding"
	"log"
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

	//mfs, err := db.GetMediaPerMimetype(m.MimeType.MIME.Value)
	//
	//for x := range mfs {
	//	log.Printf("%v",mfs[x])
	//}

	if err != nil {
		t.Error(err)
	}
}

func TestDbStorage_AddDeleteMedia(t *testing.T) {
	type fields struct {
		dbPath string
	}
	type args struct {
		media_add *adding.Media
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test db entries",
			fields:  fields{
				dbPath: testDbPath,
			},
			args:    args{media_add: &adding.Media{
				FileID:   uuid.NewV4(),
				ID:       134,
				Key:      "",
				Size:     123214,
				Filename: "testas.jpg",
				CheckSum: "57787678687",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := storage.NewDbStorage(tt.fields.dbPath)
			media, err := s.AddMedia(tt.args.media_add);
			if  (err != nil) != tt.wantErr {
				t.Errorf("DeleteMedia() error = %v, wantErr %v", err, tt.wantErr)
			}
			m2, err := s.GetByID(media.ID)
			if m2.CheckSum != media.CheckSum {
				t.Errorf("Adding media error, added = %v, found = %v", media, m2)
			}
			media_delete := deleting.Media{
				FileID: media.FileID,
				Key:    media.Key,
			}
			err = s.DeleteMedia(&media_delete);
			if  (err != nil) != tt.wantErr {
				t.Errorf("DeleteMedia() error = %v, wantErr %v", err, tt.wantErr)
			}
			m3, err := s.GetByID(media.ID)
			if m3.CheckSum != media.CheckSum {
				t.Errorf("Adding media error, added = %v, found = %v", media, m2)
			}

 		})
	}
}