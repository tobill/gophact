package viewing

import (
	"gopkg.in/h2non/filetype.v1/types"
	"github.com/satori/go.uuid"
)

//Media datamodel for adding 
type Media struct {
	FileID   uuid.UUID
	ID       uint64
	Key      string
	Size     uint64
	Filename string
	Mimetype types.Type
	CheckSum string
}

//NotFoundError error if media is not found
type NotFoundError struct {
	error string
	key string
}
