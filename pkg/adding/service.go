package adding

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"io"
	"mime/multipart"
)

//Service provides media adding
type Service interface {
	AddMedia(mpf *io.Reader, mph *multipart.FileHeader) (uint64, error)
}

//Repository access media repository
type Repository interface {
	AddMedia(media *Media) (*Media, error)
}

//FileRepository access file repository
type FileRepository interface {
	AddFile(source *io.Reader, media *Media) (error)
	DeleteDuplicateFile(media *Media) error
}

//FileRepository access file repository
type IndexRepository interface {
    FindDocuments(searchWord string) ([]string, error)
	AddDocument(media *Media) (error)
}

type service struct {
	mR Repository
	mFR FileRepository
	iR IndexRepository
}

func (e *DuplicateFileError) Error() (string) {
	return fmt.Sprintf("Duplicate File %s %s", e.key, e.error)
}

//NewService creates service
func NewService(r Repository, fr FileRepository, ir IndexRepository) Service {
	return &service{r, fr, ir}
}

//AddMedia adds media data and file
func (s *service) AddMedia(mpf *io.Reader, mph *multipart.FileHeader) (uint64, error) {
	mf := Media{}
	mf.Filename = mph.Filename
	mf.Size = uint64(mph.Size)
	mf.FileID = uuid.NewV4()

	h := sha1.New()
	var tee io.Reader
	tee = io.TeeReader(*mpf, h)


	err := s.mFR.AddFile(&tee, &mf)

	mf.CheckSum = hex.EncodeToString(h.Sum(nil))
	if err != nil { return 0, err }

	docs, err := s.iR.FindDocuments(mf.CheckSum)

	if err != nil {
		return  0, err
	}

	if len(docs) > 0 {
		err = s.mFR.DeleteDuplicateFile(&mf)
		return 0, &DuplicateFileError{error:"File already in Database",key: fmt.Sprintf("%s SHA: %v", &mf.FileID, mf.CheckSum)}
	} else {
		_, err = s.mR.AddMedia(&mf)
		s.iR.AddDocument(&mf)
	}


	return mf.ID, err
}