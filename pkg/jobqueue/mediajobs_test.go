package jobqueue_test

import (
	"gophoact/pkg/storage"
	"gophoact/pkg/editing"
	"gophoact/pkg/jobqueue"
	"testing"
	"log"
)

const testDbPath = "../../testdb/testbolt.db"
const testFilepath = "../../testdata"
const testIndexPath = "../../testdbindex"

func TestDetectMimetype(t *testing.T) {
  	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	is, err := storage.NewIndexStorage(testIndexPath) 
	defer is.CloseIndex()
	if err != nil { t.Fatal(err)	}
	e := editing.NewService(s, fs, is)
	js := jobqueue.NewService(e)

	mtjob := jobqueue.NewMimetypeJob(2, e, js) 
	mtjob.Execute()

	m, err := e.LoadMedia(2)
	if err != nil { t.Fatal(err)}
	log.Printf("%v", m)
}

func TestResizeImage(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
  defer s.CloseDb()
  if err != nil { t.Fatal(err) }
  fs := storage.NewFileStorage(testFilepath) 
  if err != nil { t.Fatal(err)	}
  is, err := storage.NewIndexStorage(testIndexPath) 
  defer is.CloseIndex()
  if err != nil { t.Fatal(err)	}
  e := editing.NewService(s, fs, is)
  js := jobqueue.NewService(e)
  mtjob := jobqueue.NewImageResizeJob(1, e, js) 
  mtjob.Execute()
  m, err := e.LoadMedia(1)
  if err != nil { t.Fatal(err)}

  doc, err := is.FindDocuments("Versions:small")

  log.Printf("%v",doc)
  log.Printf("%v", m)
}