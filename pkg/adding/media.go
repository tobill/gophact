package adding

import (
	"github.com/satori/go.uuid"
)

//Media datamodel for adding 
type Media struct {
	FileID uuid.UUID
	ID uint64
	Key string
	Size uint64
	Filename string
	CheckSum string
}


//DuplicateFileError error if media is already in repo
type DuplicateFileError struct {
	error string
	key string
}



