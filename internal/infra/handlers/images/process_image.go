package images

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type HandlerImage struct{}

func New() *HandlerImage {
	return &HandlerImage{}
}

func (h *HandlerImage) ProcessImage(im []byte, processing string) ([]byte, error) {
	srcImg, err := jpeg.Decode(bytes.NewReader(im))
	if err != nil {
		return nil, err
	}

	var dstImg *image.NRGBA

	switch processing {
	case "resize":
		dstImg = imaging.Resize(srcImg, 800, 600, imaging.Lanczos)

	case "thumbnail":
		dstImg = imaging.Thumbnail(srcImg, 100, 100, imaging.Lanczos)

	case "watermark":
		dstImg = imaging.Clone(srcImg)

		watermark := imaging.New(120, 50, color.NRGBA{0, 0, 0, 0})
		waterColor := imaging.New(120, 50, image.White)
		draw.Draw(watermark, watermark.Bounds(), waterColor, image.Point{}, draw.Over)

		dstImg = imaging.Overlay(dstImg, watermark, image.Pt(dstImg.Bounds().Dx()-120, dstImg.Bounds().Dy()-50), 0.5)
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, dstImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
