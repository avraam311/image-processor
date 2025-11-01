package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"sync"

	myKafka "github.com/avraam311/image-processor/internal/infra/kafka"
	myMinio "github.com/avraam311/image-processor/internal/infra/minio"
	"github.com/avraam311/image-processor/internal/models"
	"github.com/avraam311/image-processor/internal/repository/images"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"

	"github.com/minio/minio-go"
	"github.com/segmentio/kafka-go"
)

const (
	imageStatusProcessed = "processed"
	imageFormat          = "image/jpeg"
)

type Handler interface {
	ProcessImage([]byte, string) ([]byte, error)
}

type Repository interface {
	ChangeImageStatus(context.Context, uint, string) error
	CheckImage(context.Context, uint) error
}

type Worker struct {
	cons   *myKafka.Kafka
	cfg    *config.Config
	s3     *myMinio.Minio
	handIm Handler
	repo   Repository
}

func New(cons *myKafka.Kafka, cfg *config.Config, s3 *myMinio.Minio, handIm Handler, repo Repository) *Worker {
	return &Worker{
		cons:   cons,
		cfg:    cfg,
		s3:     s3,
		handIm: handIm,
		repo:   repo,
	}
}

func (w *Worker) Run(ctx context.Context) {
	consChan := make(chan kafka.Message)
	retryStrategy := retry.Strategy{
		Attempts: w.cfg.GetInt("retry.attempts"),
		Delay:    w.cfg.GetDuration("retry.delay"),
		Backoff:  w.cfg.GetFloat64("retry.backoff"),
	}

	go func() {
		w.cons.Consume(ctx, consChan, retryStrategy)
	}()

	var wg sync.WaitGroup
	workerCount := w.cfg.GetInt("worker.count")
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-consChan:
					imageID, err := strconv.Atoi(string(msg.Key))
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to convert msg.Key into int")
						continue
					}

					err = w.repo.CheckImage(ctx, uint(imageID))
					if err != nil {
						if errors.Is(err, images.ErrImageNotFound) {
							err := w.s3.Minio.RemoveObject(w.cfg.GetString("s3.bucket_name"), string(msg.Key))
							if err != nil {
								zlog.Logger.Warn().Err(err).Msg("worker.go - no image to process, failed to remove image from s3")
								continue
							}

							zlog.Logger.Warn().Err(err).Msg("worker.go - no image to process")
							continue
						}
					}

					imProc := models.ImageKafka{}
					err = json.Unmarshal(msg.Value, &imProc)
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to unmarshal message into struct")
						continue
					}

					s3Object, err := w.s3.Minio.GetObject(w.cfg.GetString("s3.bucket_name"), string(msg.Key), minio.GetObjectOptions{})
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to get image from s3")
						continue
					}
					defer s3Object.Close()
					buf := new(bytes.Buffer)
					if _, err = io.Copy(buf, s3Object); err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to copy s3Object into buffer")
						continue
					}
					imageBytes := buf.Bytes()
					image := models.Image{}
					err = json.Unmarshal(imageBytes, &image)
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to unmarshal buf image into struct")
						continue
					}

					processedImage, err := w.handIm.ProcessImage(image.Image, imProc.Processing)
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to process image")
						continue
					}

					objectName := string(msg.Key)
					imageAsReader := bytes.NewReader(processedImage)
					size := int64(len(processedImage))
					putObjectOptions := minio.PutObjectOptions{
						ContentType: imageFormat,
					}
					_, err = w.s3.Minio.PutObject(w.cfg.GetString("s3.bucket_name"), objectName, imageAsReader, size, putObjectOptions)
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to put processed image into s3")
						continue
					}

					err = w.repo.ChangeImageStatus(ctx, (uint(imageID)), imageStatusProcessed)
					if err != nil {
						zlog.Logger.Warn().Err(err).Msg("worker.go - failed to change image status")
						continue
					}
					zlog.Logger.Info().Interface("image", msg).Msg("image is processed")
				}
			}
		}(i)
	}

	wg.Wait()
}
