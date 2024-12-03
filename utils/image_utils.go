package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

const asciiRamp = "@#S%?*+;:,. "

// LoadImage loads an image from a file
func LoadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Grayscale converts an image to grayscale
func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			originalColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}
	return grayImg
}

// GenerateASCIIArt generates ASCII art from a grayscale image
func GenerateASCIIArt(grayImg *image.Gray) string {
	bounds := grayImg.Bounds()
	var asciiArt strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayValue := grayImg.GrayAt(x, y).Y
			asciiArt.WriteString(pixelToASCII(grayValue))
		}
		asciiArt.WriteString("\n")
	}

	return asciiArt.String()
}

// pixelToASCII maps a grayscale value to an ASCII character
func pixelToASCII(gray uint8) string {
	index := int(gray) * (len(asciiRamp) - 1) / 255
	return string(asciiRamp[index])
}
