package jobqueue

import (
	"gophoact/pkg/editing"
	"log"
	"gopkg.in/h2non/filetype.v1"
)

type mimetypeJob struct {
	status int
	jobName string 
	objectID uint64
	edit editing.Service
}

func (j *mimetypeJob) Execute() (int) {
	log.Printf("loading objId %d", j.objectID)
	m, fob, err := j.edit.LoadMediaWithFiledata(j.objectID)
	log.Print("executing mimetype job")
	defer fob.Close()
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	b := make([]byte, 512)
	_, err = fob.Read(b)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	t, err := filetype.Match(b)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	m.MimeType = t
	j.edit.SaveMedia(m)
	return j.status
}

func (j *mimetypeJob) SetStatus(status int) {
	j.status = status
}

func (j *mimetypeJob) GetStatus() (int){
	return j.status
}

func (j *mimetypeJob) GetJobName() (string) {
	return j.jobName
}

func (j *mimetypeJob) GetObjectID() (uint64){
	return j.objectID
}

//NewMimtypeJob returns job for detecting mimetype of media
func NewMimetypeJob(objectID uint64, e editing.Service) (Job){
	mj := mimetypeJob{
		jobName: "Mimetype",
		objectID: objectID,
		status: 0,
		edit: e,
	}
	return &mj
}