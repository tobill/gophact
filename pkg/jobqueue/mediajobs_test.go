package jobqueue_test

import (
	"gophoact/pkg/editing"
	"gophoact/pkg/jobqueue"
	"gophoact/pkg/storage"
	"log"
	"testing"
)

const testDbPath = "../../testdb/testbolt.db"
const testFilepath = "../../testdata"
const testIndexPath = "../../testdbindex"

func TestDetectMimetype(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	if err != nil {
		t.Fatal(err)
	}
	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()
	if err != nil {
		t.Fatal(err)
	}
	e := editing.NewService(s, fs, is)

	mtjob := jobqueue.NewMimetypeJob(2, e)
	mtjob.Execute()

	m, err := e.LoadMedia(2)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("%v", m)
}

func TestResizeImage(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	if err != nil {
		t.Fatal(err)
	}
	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()
	if err != nil {
		t.Fatal(err)
	}
	e := editing.NewService(s, fs, is)
	mtjob := jobqueue.NewImageResizeJob(1, e)
	mtjob.Execute()
	m, err := e.LoadMedia(1)
	if err != nil {
		t.Fatal(err)
	}

	doc, err := is.FindDocuments("Versions:small")

	log.Printf("%v", doc)
	log.Printf("%v", m)
}

func TestCheckSum(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil {
		t.Fatal(err)
	}
	fs := storage.NewFileStorage(testFilepath)
	if err != nil {
		t.Fatal(err)
	}
	is, err := storage.NewIndexStorage(testIndexPath)
	defer is.CloseIndex()
	if err != nil {
		t.Fatal(err)
	}
	e := editing.NewService(s, fs, is)
	mtjob := jobqueue.NewChkSumJob(1, e)
	mtjob.Execute()
	m, err := e.LoadMedia(1)
	if err != nil {
		t.Fatal(err)
	}
	chksum := "39a641f0d33003a5a729a8c2415f4b0b130e4620"

	if chksum != m.CheckSum {
		t.Fatal("chk wrong chcksum")
	}
}
