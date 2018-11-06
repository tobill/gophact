package editing

import (
	"os"
)

//Service provides media editing
type Service interface {
	LoadMedia(objectID uint64) (*Media, error)
	LoadMediaWithFiledata(objectID uint64) (*Media, *os.File, error)
	SaveMedia(*Media) error
}

//Repository access media repository
type Repository interface {
	LoadMedia(objectID uint64) (*Media, error)
	SaveMedia(*Media) error
}

//FileRepository access file repository
type FileRepository interface {
	GetFile(fileID string) (*os.File, error)
}

type service struct {
	mR Repository
	mFR FileRepository
}

//NewService creates service
func NewService(r Repository, fr FileRepository) Service {
	return &service{r, fr}
}


func (s *service) LoadMedia(objectID uint64) (*Media, error) {
	return s.mR.LoadMedia(objectID)
}

func (s *service) LoadMediaWithFiledata(objectID uint64) (*Media, *os.File, error) {
	m, err := s.mR.LoadMedia(objectID)
	if err != nil {
		return m, nil, err
	}
	fob, err := s.mFR.GetFile(m.FileID.String())
	return m, fob, err
}

func (s *service) SaveMedia(m *Media) (error) {
	return s.mR.SaveMedia(m)
}