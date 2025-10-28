package images

import (
	"context"

	"github.com/avraam311/image-processor/internal/models"

	"github.com/go-playground/validator/v10"
)

type Service interface {
	UploadImage(context.Context, *models.Image) (uint, error)
	GetProcessedImage(context.Context, uint) ([]byte, error)
	DeleteImage(context.Context, uint) error
}

type Handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(service Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}
