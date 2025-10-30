package images

import (
	"context"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/kafka"

	"github.com/minio/minio-go"
)

type Repository interface {
	UploadImage(context.Context, string) (uint, error)
	GetProcessedImage(context.Context, uint) error
	ChangeImageStatus(context.Context, uint) error
}

type Service struct {
	repo Repository
	prod *kafka.Producer
	cfg  *config.Config
	s3   *minio.Client
}

func NewService(repo Repository, prod *kafka.Producer, cfg *config.Config, s3 *minio.Client) *Service {
	return &Service{
		repo: repo,
		prod: prod,
		cfg:  cfg,
		s3:   s3,
	}
}
