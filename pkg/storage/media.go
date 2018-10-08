package storage

import (
	"github.com/satori/go.uuid"
)

//Media defines the storage form for media objects
type Media struct {
	FileID uuid.UUID
	ID uint64
	Key string
	Size uint64
	Filename string
}