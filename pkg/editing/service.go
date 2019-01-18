package editing

import (
	"io"
	"os"
)

//Service provides media editing
type Service interface {
	LoadMedia(objectID uint64) (*Media, error)
	LoadMediaWithFiledata(objectID uint64) (*Media, *os.File, error)
	AddSmallImage(objectID uint64, mr *io.Reader) (error)
	SaveMedia(*Media) error
}

//Repository access media repository
type Repository interface {
	LoadMedia(objectID uint64) (*Media, error)
	SaveMedia(*Media) error
}

//FileRepository access file repository
type FileRepository interface {
	GetOriginalFile(fileID string) (*os.File, error)
	AddSmallFile(source *io.Reader, media *Media) error
}

//IndexRepository access index / search repository
type IndexRepository interface {
	AddDocument(m *Media) (error)
}


type service struct {
	mR Repository
	mFR FileRepository
	mIR  IndexRepository
}

//NewService creates service
func NewService(r Repository, fr FileRepository, ir IndexRepository) Service {
	return &service{r, fr, ir}
}

func (s * service) OpenSmallImageWriter(objectID uint64) (*io.Writer, error) {
	return nil, nil
}	

func (s *service) LoadMedia(objectID uint64) (*Media, error) {
	return s.mR.LoadMedia(objectID)
}

func (s *service) LoadMediaWithFiledata(objectID uint64) (*Media, *os.File, error) {
	m, err := s.mR.LoadMedia(objectID)
	if err != nil {
		return m, nil, err
	}
	fob, err := s.mFR.GetOriginalFile(m.FileID.String())
	return m, fob, err
}

func (s *service) SaveMedia(m *Media) (error) {
	err := s.mR.SaveMedia(m)
	if err != nil {
		return err
	}
	return s.mIR.AddDocument(m)


}

func (s *service) AddSmallImage(objectID uint64, mr *io.Reader) (error) {
	m, err := s.mR.LoadMedia(objectID)
	if err != nil { return err }
	err = s.mFR.AddSmallFile(mr, m)
	if err != nil { return err }
	return err
}