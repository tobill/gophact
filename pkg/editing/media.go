package editing

import (
	"github.com/satori/go.uuid"
	"gopkg.in/h2non/filetype.v1/types"
)

//Media datamodel for adding 
type Media struct {
	FileID uuid.UUID
	ID uint64
	Key string
	Size uint64
	Filename string
	MimeType types.Type
	Versions []string
}

//NotFoundError error if media is not found
type NotFoundError struct {
	error string
	key string
}
