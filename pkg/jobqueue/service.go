package jobqueue

import (
	"io"
	"log"
	"gophoact/pkg/editing"
    "sync"
    "os"
)


var  jobqueue chan Job
var wg sync.WaitGroup

//Service for jobs
type Service interface {
    EnqueueJob(Job)
    CloseQueue()
    GenerateMimetypeJob(objetcID uint64)
    GenerateResizeJob(objetcID uint64)
}

//EditingService dummy
type EditingService interface {
    LoadMedia(objectID uint64) (*editing.Media, error)
    LoadMediaWithFiledata(objectID uint64) (*editing.Media, *os.File, error)
    AddSmallImage(objectID uint64, mr *io.Reader) (error)
	SaveMedia(*editing.Media) error
}

type service struct {
    queue chan Job
    es EditingService
}

//NewService generate new service holding queue
func NewService(es EditingService) (Service) {
    wg.Add(1)
    jq  := make(chan Job, 100)
    go worker(jq)
    return &service{jq, es}	
}

func worker(jobChan <-chan Job) {
    defer wg.Done()
    for job := range jobChan {
        log.Printf("doing the job")
        job.Execute()
        job.SetStatus(1)
    }
}

func (s *service) EnqueueJob(j Job) {
    log.Printf("enqueue")
    s.queue <- j
}

func (s *service) GenerateMimetypeJob(objetcID uint64) {
    log.Printf("generate")
    mtjob := NewMimetypeJob(objetcID, s.es) 
    s.EnqueueJob(mtjob)
}

func (s *service) GenerateResizeJob(objetcID uint64) {
    log.Printf("generate")
    mtjob := NewImageResizeJob(objetcID, s.es) 
    s.EnqueueJob(mtjob)
}



func (s *service) CloseQueue () {
    close(s.queue)
    wg.Wait()
} 

