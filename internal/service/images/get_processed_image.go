package images

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/minio/minio-go"
)

func (s *Service) GetProcessedImage(ctx context.Context, id uint) ([]byte, error) {
	err := s.repo.CheckImage(ctx, id)
	if err != nil {
		return []byte{}, fmt.Errorf("service/images - %w", err)
	}

	getObjectOptions := minio.GetObjectOptions{}
	objectName := strconv.Itoa(int(id))
	object, err := s.s3.GetObject(s.cfg.GetString("s3.bucket_name"), objectName, getObjectOptions)
	if err != nil {
		return []byte{}, fmt.Errorf("service/images - failed to get image from s3 - %w", err)
	}
	defer object.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, object); err != nil {
		return nil, fmt.Errorf("service/images - failed to copy object into buffer - %w", err)
	}
	image := buf.Bytes()

	return image, nil
}
