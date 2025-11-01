package models

type Image struct {
	Image      []byte `json:"image" validate:"required"`
	Processing string `json:"processing" validate:"required"`
}

type ImageKafka struct {
	Processing string `json:"processing" validate:"required"`
}
