package images

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/avraam311/image-processor/internal/models"

	"github.com/wb-go/wbf/retry"

	"github.com/minio/minio-go"
)

const (
	imageStatusInProcess = "in process"
	imageFormat          = "image/jpeg"
)

func (s *Service) UploadImage(ctx context.Context, im *models.Image) (uint, error) {
	id, err := s.repo.UploadImage(ctx, imageStatusInProcess)
	if err != nil {
		return 0, fmt.Errorf("service/upload_image.go - %w", err)
	}

	prodStrategy := retry.Strategy{
		Attempts: s.cfg.GetInt("retry.attempts"),
		Delay:    s.cfg.GetDuration("retry.delay"),
		Backoff:  s.cfg.GetFloat64("retry.backoff"),
	}
	imageKey, err := json.Marshal(id)
	if err != nil {
		return 0, fmt.Errorf("service/upload_image.go - failed to marshal id into json - %w", err)
	}
	imageValue, err := json.Marshal(im.Processing)
	if err != nil {
		return 0, fmt.Errorf("service/upload_image.go - failed to marshal processing into json - %w", err)
	}
	err = s.prod.SendWithRetry(ctx, prodStrategy, imageKey, imageValue)
	if err != nil {
		return 0, fmt.Errorf("service/upload_image.go - failed to send request to kafka - %w", err)
	}

	objectName := strconv.Itoa(int(id))
	imageAsReader := bytes.NewReader(imageValue)
	size := int64(len(imageValue))
	putObjectOptions := minio.PutObjectOptions{
		ContentType: imageFormat,
	}
	_, err = s.s3.PutObject(s.cfg.GetString("s3.bucket_name"), objectName, imageAsReader, size, putObjectOptions)
	if err != nil {
		return 0, fmt.Errorf("service/upload_image.go - failed to put image in s3 - %w", err)
	}

	return id, nil
}
