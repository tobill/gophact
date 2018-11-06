package jobqueue_test

import (
	"gophoact/pkg/editing"
	"gophoact/pkg/storage"
	"log"
	"testing"
	"gophoact/pkg/jobqueue"
)

type emptyJob struct {
	status int
	jobName string 
	objectId uint64
} 

func (j *emptyJob) Execute() (int) {
	log.Printf("doing nothing but executing")
	return j.status
}

func (j *emptyJob) SetStatus(status int) {
	j.status = status
}

func (j *emptyJob) GetStatus() (int){
	return j.status
}

func (j *emptyJob) GetJobName() (string) {
	return j.jobName
}

func (j *emptyJob) GetObjectID() (uint64){
	return j.objectId
}

func NewEmptyJob() (jobqueue.Job) {
	return &emptyJob{
		status: 1,
		jobName: "EmptyJob",
		objectId: 0,
		} 
}

func TestExecuteJob(t *testing.T) {
	s, err := storage.NewDbStorage(testDbPath)
	defer s.CloseDb()
	if err != nil { t.Fatal(err) }
	fs := storage.NewFileStorage(testFilepath) 
	if err != nil { t.Fatal(err)	}
	e := editing.NewService(s, fs)

	job := emptyJob {}
	service := jobqueue.NewService(e)
	defer service.CloseQueue()
	service.EnqueueJob(&job)
	service.EnqueueJob(&job)
	service.EnqueueJob(&job)
	service.EnqueueJob(&job)
	service.EnqueueJob(&job)
	service.EnqueueJob(&job)
	log.Printf("finisehd enque")
}
