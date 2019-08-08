package viewing

import (
	"fmt"
	"os"
)

//Service provides media adding
type Service interface {
	ListAll(numStart uint64, numItems uint64) ([]*Media, error)
	GetByID(id uint64) (*Media, error)
	GetFileByID(id uint64) (*os.File, error)
}

//Repository access media repository
type Repository interface {
	ListAll(start uint64, numOfItems uint64) ([]*Media, error)
	GetByID(id uint64) (*Media, error)
}

//FileRepository access media repository
type FileRepository interface {
	GetOriginalFile(fileID string) (*os.File, error) 
}

type service struct {
	mR Repository
	mFR FileRepository
}

//NewService creates service
func NewService(r Repository, fr FileRepository) Service {
	return &service{r, fr}
}

func (e *NotFoundError) Error() (string) {
	return fmt.Sprintf("Key \"%s\" Not Found: %s", e.key, e.error)
}

func (s *service) ListAll(numStart uint64, numItems uint64) ([]*Media, error) {
	return s.mR.ListAll(numStart, numItems)
}

func (s *service) GetByID(id uint64) (*Media, error) {
	return s.mR.GetByID(id)
}

func (s *service) GetFileByID(id uint64) (*os.File, error) {
	m, err := s.mR.GetByID(id)
	if err != nil {
		return nil, &NotFoundError{error:"Not found in Database",key: fmt.Sprintf("%d", id)}
	}
	fop, err := s.mFR.GetOriginalFile(m.FileID.String())
	if err != nil {
		return nil, &NotFoundError{error:"Not found in Filesystem",key: fmt.Sprintf("%s", m.FileID)}
	}
	return fop, nil
}

