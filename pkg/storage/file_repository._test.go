package storage_test

import (
	uuid "github.com/satori/go.uuid"
	"gophoact/pkg/adding"
	"gophoact/pkg/deleting"
	"gophoact/pkg/storage"
	"io"
	"os"
	"testing"
)

func TestFileStorage_AddDelete(t *testing.T) {
	fp, err := os.Open(testFile)
	defer fp.Close()
	if err != nil {
		t.Error(err)
	}
	type fields struct {
		dirpathRoot     string
	}
	type args struct {
		media_del deleting.Media
		media_add adding.Media
		fp *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test1",
			fields:  fields{
				dirpathRoot: testFilepath,
			},
			args:    args{
				media_add: adding.Media{
					Filename: "Tesddddtfile",
					Size: 12345,
					FileID: uuid.NewV4(),
				},
				fp: fp,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewFileStorage(tt.fields.dirpathRoot)
			var fpr io.Reader = tt.args.fp
			if err := s.AddFile(&fpr, &tt.args.media_add); err !=nil {
				t.Errorf("DeleteMedia, want to add  error = %v", err)
			}
			fp2, err := s.GetOriginalFile(tt.args.media_add.FileID.String())
			if err != nil {
				t.Errorf("Added not found,  error = %v", err)
			}
			err = fp2.Close();
			if  err != nil {
				t.Errorf("Added not found,  error = %v", err)
			}
			media_del :=  deleting.Media{
				FileID: tt.args.media_add.FileID,
				Key:    tt.args.media_add.Key,
			}
			if err := s.DeleteMedia(&media_del); (err != nil) != tt.wantErr {
				t.Errorf("DeleteMedia() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
