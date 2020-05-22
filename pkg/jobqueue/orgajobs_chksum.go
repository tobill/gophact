package jobqueue

import (
	"crypto/sha1"
	"encoding/hex"
	"gophoact/pkg/editing"
	"io"
	"log"
)

type chkSumJob struct {
	status int
	jobName string 
	objectID uint64
	edit editing.Service
}

func (j *chkSumJob) Execute() (int) {
	log.Printf("loading objId %d", j.objectID)
	m, fob, err := j.edit.LoadMediaWithFiledata(j.objectID)
	log.Print("executing chksum job")
	defer fob.Close()
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	h := sha1.New()
	_, err = io.Copy(h, fob)
	if err != nil {
		log.Printf("error doing my job %v", err)
		j.status = 3
		return 3
	}
	m.CheckSum = hex.EncodeToString(h.Sum(nil))
	j.edit.SaveMedia(m)
	return j.status
}

func (j *chkSumJob) SetStatus(status int) {
	j.status = status
}

func (j *chkSumJob) GetStatus() (int){
	return j.status
}

func (j *chkSumJob) GetJobName() (string) {
	return j.jobName
}

func (j *chkSumJob) GetObjectID() (uint64){
	return j.objectID
}

//NewImageResizeJob returns job for detecting mimetype of media
func NewChkSumJob(objectID uint64, e editing.Service) (Job){
	mj := chkSumJob{
		jobName: "ResizeImage",
		objectID: objectID,
		status: 0,
		edit: e,
	}
	return &mj
}