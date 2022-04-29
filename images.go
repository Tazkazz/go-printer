package go_printer

import (
	"errors"
	"github.com/lestrrat-go/dither"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"math"
)

const (
	maxWidthPixels  int     = 384
	errorMultiplier float32 = 1.18
)

var thumbnailFilter = resize.Lanczos3
var monochromeFilter = dither.Stucki

type PrintableImage struct {
	lines         *[]*[]uint8
	width, height int
}

func NewPrintableImage(
	srcImg image.Image,
	alignment AlignmentLevel,
) (img *PrintableImage, err error) {
	if srcImg == nil {
		err = errors.New("srcImg cannot be nil")
		return
	}

	resizedImg := resizeImageIfTooWide(srcImg)
	filledImage := fillImageAlphaAndAlign(resizedImg, alignment)
	monoImage := convertToMonoImage(filledImage)
	lines, width, height := splitImageToLines(monoImage)

	img = &PrintableImage{
		lines:  &lines,
		width:  width,
		height: height,
	}
	return
}

func resizeImageIfTooWide(img image.Image) image.Image {
	return resize.Thumbnail(uint(maxWidthPixels), math.MaxUint, img, thumbnailFilter)
}

func fillImageAlphaAndAlign(img image.Image, alignment AlignmentLevel) image.Image {
	newImg := image.NewRGBA(image.Rectangle{Max: image.Point{X: maxWidthPixels, Y: img.Bounds().Dy()}})
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: image.White}, image.Point{}, draw.Src)

	x := 0
	if img.Bounds().Dx() < maxWidthPixels {
		switch alignment {
		case AlignmentCenter:
			x = int(math.Round(float64(maxWidthPixels-img.Bounds().Dx()) / 2.0))
		case AlignmentRight:
			x = maxWidthPixels - img.Bounds().Dx()
		}
	}

	draw.Draw(newImg, newImg.Bounds(), img, image.Point{X: -x}, draw.Over)

	return newImg
}

func convertToMonoImage(img image.Image) *image.Gray {
	return dither.Monochrome(monochromeFilter, img, errorMultiplier).(*image.Gray)
}

func splitImageToLines(img *image.Gray) (lines []*[]uint8, width, height int) {
	realWidth := img.Bounds().Dx()
	offset := (8 - realWidth%8) % 8
	width = realWidth / 8

	height = img.Bounds().Dy()

	if offset != 0 {
		width++
	}

	lines = make([]*[]uint8, 0, height)

	for y := 0; y < height; y++ {
		line := make([]uint8, 0, width)
		bits := uint8(0)

		for x := 0; x < realWidth; x++ {
			bits <<= 1
			if img.GrayAt(x, y).Y < 128 {
				bits += 1
			}
			if (x+1)%8 == 0 {
				line = append(line, bits)
				bits = 0
			}
		}

		if offset != 0 {
			bits <<= offset
			line = append(line, bits)
		}

		lines = append(lines, &line)
	}

	return
}
