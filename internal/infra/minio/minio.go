package minio

import (
	"github.com/minio/minio-go"
)

type Minio struct {
	Minio *minio.Client
}

func New(endpoint, user, password, bucketName, location string, ssl bool) (*Minio, error) {
	minioClient, err := minio.New(endpoint, user, password, ssl)
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = minioClient.MakeBucket(bucketName, location)
		if err != nil {
			return nil, err
		}
	}

	return &Minio{
		Minio: minioClient,
	}, nil
}
