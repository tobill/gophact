package storage

import (
	"path/filepath"
	"os"
	"io"
	"gophoact/pkg/adding"
	"gophoact/pkg/editing"
)

//FileStorage stores filedata on disk 
type FileStorage struct {
	dirpathRoot string
	dirpathOriginal string
	dirpathSmall string
	dirpathThumb string

}


//NewFileStorage create new storage for files
func NewFileStorage(dirpath string) (*FileStorage) {
	s := FileStorage{
		dirpathRoot: dirpath,
		dirpathOriginal: "original",
		dirpathSmall:  "small",
		dirpathThumb: "thumb",
	}
	return &s
}

func (s *FileStorage) getFilePath(fileID string, size string) (string, error) {
	subdir := filepath.Join(fileID[0:1], fileID[1:2])
	fpath := filepath.Join(s.dirpathRoot, size, subdir, fileID)
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return "", err
	}
	return fpath, nil
}

func (s *FileStorage) checkSubDirs(fileID string, size string) (string, error) {
	subdir := filepath.Join(fileID[0:1], fileID[1:2])
	if _, err := os.Stat(filepath.Join(s.dirpathRoot, size, subdir)); os.IsNotExist(err) {
		if _, err := os.Stat(filepath.Join(s.dirpathRoot, size, fileID[0:1])); os.IsNotExist(err) {
			err = os.Mkdir(filepath.Join(s.dirpathRoot, size, fileID[0:1]), os.ModeDir)
			if err != nil { return "", err}
		}
		err = os.Mkdir(filepath.Join(s.dirpathRoot, size, subdir), os.ModeDir)
		if err != nil { return "", err}
	}
	return filepath.Join(subdir,fileID), nil
}

//GetOriginalFile returns origina file by uuid
func (s *FileStorage) GetOriginalFile(fileID string) (*os.File, error) {
	fpath, err := s.getFilePath(fileID, s.dirpathOriginal)
	if (err != nil ) { return nil, err }
	return os.Open(fpath)
}

//GetSmallFile
func (s *FileStorage) GetSmallFile(fileID string) (*os.File, error) {
	fpath, err := s.getFilePath(fileID, s.dirpathSmall)
	if (err != nil ) { return nil, err }
	return os.Open(fpath)
}

//AddFile creates file in repo
func (s *FileStorage) AddFile(source *io.Reader, media *adding.Media) error {
	filep, err := s.checkSubDirs(media.FileID.String(), s.dirpathOriginal)
	if err != nil { return err }
	fpath := filepath.Join(s.dirpathRoot, s.dirpathOriginal, filep)
	fd, err := os.Create(fpath)
	if err != nil { return err}
	defer fd.Close()
	_, err = io.Copy(fd, *source)
	return err
}

//AddSmallFile creates file in repo
func (s *FileStorage) AddSmallFile(source *io.Reader, media *editing.Media) error {
	filep, err := s.checkSubDirs(media.FileID.String(), s.dirpathSmall)
	if err != nil { return err }
	fpath := filepath.Join(s.dirpathRoot, s.dirpathSmall, filep)
	fd, err := os.Create(fpath)
	if err != nil { return err}
	defer fd.Close()
	_, err = io.Copy(fd, *source)
	return err
}

