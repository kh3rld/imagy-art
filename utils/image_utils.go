package utils

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"sync"
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

// RenderASCIIToImage renders ASCII art as an image and saves it as a PNG
func RenderASCIIToImage(asciiArt string, file string, fontSize int) error {
	lines := strings.Split(asciiArt, "\n")

	for len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return fmt.Errorf("ASCII art is empty after trimming")
	}

	width := len(lines[0])
	height := len(lines)

	img := image.NewGray(image.Rect(0, 0, width*fontSize, height*fontSize))
	var wg sync.WaitGroup
	for y := 0; y < height; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			for x := 0; x < width; x++ {
				asciiChar := lines[y][x]
				grayValue := (len(asciiRamp) - 1 - asciiRampIndex(asciiChar)) * 255 / (len(asciiRamp) - 1)
				col := color.Gray{Y: uint8(grayValue)}
				for i := 0; i < fontSize; i++ {
					for j := 0; j < fontSize; j++ {
						img.Set(x*fontSize+i, y*fontSize+j, col)
					}
				}
			}
		}(y)
	}
	wg.Wait()

	outFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	return png.Encode(writer, img)
}

// asciiRampIndex returns the index of an ASCII character in the ASCII ramp
func asciiRampIndex(char byte) int {
	for i, c := range asciiRamp {
		if c == rune(char) {
			return i
		}
	}
	return len(asciiRamp) - 1 // Default to lightest if character not found
}
