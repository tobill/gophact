package adding

import (
	"mime/multipart"
	"github.com/satori/go.uuid"
	"io"
)

//Service provides media adding
type Service interface {
	AddMedia(mpf *io.Reader, mph *multipart.FileHeader) error

}

//Repository access media repository
type Repository interface {
	AddMedia(media *Media) error
}

//FileRepository access file repository
type FileRepository interface {
	AddFile(source *io.Reader, media *Media) error
}


type service struct {
	mR Repository
	mFR FileRepository
}

//NewService creates service
func NewService(r Repository, fr FileRepository) Service {
	return &service{r, fr}
}

//AddMedia adds media data and file
func (s *service) AddMedia(mpf *io.Reader, mph *multipart.FileHeader) error {
	mf := Media{}
	mf.Filename = mph.Filename
	mf.Size = uint64(mph.Size)
	mf.FileID = uuid.NewV4()
	err := s.mFR.AddFile(mpf, &mf)
	if err != nil { return err }
	s.mR.AddMedia(&mf)
	return err
}