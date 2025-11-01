package images

import (
	"context"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/kafka"

	"github.com/avraam311/image-processor/internal/infra/minio"
)

type Repository interface {
	SetImageStatus(context.Context, string) (uint, error)
	CheckImage(context.Context, uint) error
	DeleteImage(context.Context, uint) error
}

type Service struct {
	repo Repository
	prod *kafka.Producer
	cfg  *config.Config
	s3   *minio.Minio
}

func NewService(repo Repository, prod *kafka.Producer, cfg *config.Config, s3 *minio.Minio) *Service {
	return &Service{
		repo: repo,
		prod: prod,
		cfg:  cfg,
		s3:   s3,
	}
}
