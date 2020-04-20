package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Note type used for defining a note
type Note struct {
	ID      uuid.UUID
	Content string
	Date    time.Time
}

var (
	stored []*Note
)

func init() {

}

// GetByID retrieves a note based on it's ID
func GetByID(id string) (*Note, error) {
	_, note, err := searchCollection(id)

	return note, err
}

//Store retrieves a note based on it's ID
func Store(content string) (Note, error) {
	note := Note{
		ID:      uuid.New(),
		Content: content,
		Date:    time.Now(),
	}
	stored = append(stored, &note)

	return note, nil
}

// Fetch is used to retrieve a note and remove it from the collection at the same time
func Fetch(id string) (*Note, error) {
	index, note, err := searchCollection(id)

	if -1 != index && nil == err {
		stored = append(stored[:index], stored[index+1:]...)
	}

	return note, err
}

func searchCollection(id string) (int, *Note, error) {
	for index, item := range stored {
		if item.ID.String() == id {
			return index, item, nil
		}
	}

	return -1, &Note{}, errors.New("note does not exist")
}
