package deleting

//Service provides media deleting
type Service interface {
	DeleteMedia(media *Media) error
}

type IndexRepository interface {
	DeleteMedia(media *Media) error
}

type Repository interface {
	DeleteMedia(media *Media) error
}

type FileRepository interface {
	DeleteMedia(media *Media) error
}


type service struct {
	mR Repository
	mFR FileRepository
	mIR IndexRepository
}

//NewService creates new Service
func NewService(r Repository, fR FileRepository, iR IndexRepository) Service {
	return &service{r, fR, iR}
}

func (s *service) DeleteMedia(media *Media) error {
	err := s.mFR.DeleteMedia(media)
	if err != nil {
		return err
	}
	err = s.mIR.DeleteMedia(media)
	if err != nil {
		return err
	}
	err = s.mR.DeleteMedia(media)
	return err
}