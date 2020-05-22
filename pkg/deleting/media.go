package deleting

import uuid "github.com/satori/go.uuid"

//Media datamodel for deleting
type Media struct {
	FileID uuid.UUID
	Key string
}

