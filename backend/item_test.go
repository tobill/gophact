package backend

import (
	"log"
	"testing"
	
)

func TestPhotoSave(t *testing.T) {
	var p IMedia
	Connect(testPath)
	defer Close()
	p = &MediaFile{Filename: "tsest.jpg", ID: 5}
	err := p.Save()
	log.Println(p)
	if err != nil {
		t.Error(err)
	}
}


func TestPhotoLoadAll(t *testing.T) {
	var allPhotos []*MediaFile
	allPhotos, err := LoadAll()
	if err != nil {
		t.Error(err)
	}
	if len(allPhotos) == 0 {
		t.Fatal("data is empty")
	} 
}