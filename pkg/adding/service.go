package adding

import (
	"mime/multipart"
	"github.com/satori/go.uuid"
	"gophoact/pkg/jobqueue"
	"io"
)

//Service provides media adding
type Service interface {
	AddMedia(mpf *io.Reader, mph *multipart.FileHeader) (uint64, error)

}

//Repository access media repository
type Repository interface {
	AddMedia(media *Media) (uint64, error)
}

//FileRepository access file repository
type FileRepository interface {
	AddFile(source *io.Reader, media *Media) (error)
}

//Jobqueue manage jobs
type Jobqueue interface {
	EnqueueJob(j jobqueue.Job)
}

type service struct {
	mR Repository
	mFR FileRepository
	jq Jobqueue
}

//NewService creates service
func NewService(r Repository, fr FileRepository, jq Jobqueue) Service {
	return &service{r, fr, jq}
}

//AddMedia adds media data and file
func (s *service) AddMedia(mpf *io.Reader, mph *multipart.FileHeader) (uint64, error) {
	mf := Media{}
	mf.Filename = mph.Filename
	mf.Size = uint64(mph.Size)
	mf.FileID = uuid.NewV4()
	err := s.mFR.AddFile(mpf, &mf)
	if err != nil { return 0, err }
	id, err := s.mR.AddMedia(&mf)
	
	return id, err
}