package jobqueue_test

import (
	"gophoact/pkg/storage"
	"gophoact/pkg/editing"
	"gophoact/pkg/jobqueue"
	"testing"
	"log"
)

const testDbPath = "../../testdb"
const testFilepath = "../../testdata"

func TestDetectMimetype(t *testing.T) {
  	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	e := editing.NewService(s, fs)

	mtjob := jobqueue.NewMimetypeJob(1, e) 
	mtjob.Execute()

	m, err := e.LoadMedia(1)
	if err != nil { t.Fatal(err)}
	log.Printf("%v", m)
}