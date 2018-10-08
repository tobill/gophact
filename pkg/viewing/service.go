package viewing

import (
//	"github.com/satori/go.uuid"
)

//Service provides media adding
type Service interface {
	ListAll() ([]*Media, error)

}

//Repository access media repository
type Repository interface {
	ListAll() ([]*Media, error)
}

//FileRepository access media repository
type FileRepository interface {
}

type service struct {
	mR Repository
	mFR FileRepository
}

//NewService creates service
func NewService(r Repository, fr FileRepository) Service {
	return &service{r, fr}
}

func (s *service) ListAll() ([]*Media, error) {
	return s.mR.ListAll()
}