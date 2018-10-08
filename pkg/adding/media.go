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
}
