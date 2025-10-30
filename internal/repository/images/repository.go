package images

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrImageNotFound  = errors.New("image not found")
	ErrImageInProcess = errors.New("image in process")
)

type Repository struct {
	db *dbpg.DB
}

func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}
