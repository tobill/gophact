package jobqueue

import (
	"io"
	"os"
	"io/ioutil"
	"gophoact/pkg/editing"
	"log"
	"github.com/disintegration/imaging"
)

type imageResizeJob struct {
	status int
	jobName string 
	objectID uint64
	edit editing.Service
}

func (j *imageResizeJob) Execute() (int) {
	log.Printf("loading objId %d", j.objectID)
	
	
	m, fob, err := j.edit.LoadMediaWithFiledata(j.objectID)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	log.Print("executing image resize job")
	defer fob.Close()
	
	img, err := imaging.Decode(fob)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	t1 := img.Bounds()
	log.Printf("%v", t1)

	img = imaging.Resize(img, 1024, 0, imaging.Lanczos)

	t1 = img.Bounds()
	log.Printf("%v", t1)
	
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	tmpfile, err := ioutil.TempFile("", "img")
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	defer os.Remove(tmpfile.Name()) // clean up
	err = imaging.Encode(tmpfile, img, imaging.JPEG)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	tmpfile.Close()
	tmpr, err := os.Open(tmpfile.Name()) 
	defer tmpr.Close()
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	var tmpread io.Reader = tmpr
	err = j.edit.AddSmallImage(j.objectID, &tmpread)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	found := false
	for _, v := range m.Versions {
		if v == "small" {
			found = true
		}
	}
	if !found {
		m.Versions = append(m.Versions, "small")
	}
	j.edit.SaveMedia(m)
	return j.status
}

func (j *imageResizeJob) SetStatus(status int) {
	j.status = status
}

func (j *imageResizeJob) GetStatus() (int){
	return j.status
}

func (j *imageResizeJob) GetJobName() (string) {
	return j.jobName
}

func (j *imageResizeJob) GetObjectID() (uint64){
	return j.objectID
}

//NewImageResizeJob returns job for detecting mimetype of media
func NewImageResizeJob(objectID uint64, e editing.Service) (Job){
	mj := imageResizeJob{
		jobName: "ResizeImage",
		objectID: objectID,
		status: 0,
		edit: e,
	}
	return &mj
}