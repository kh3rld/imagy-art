package main

import (
	"fmt"
	"log"

	"github.com/kh3rld/imagy-art/utils"
)

const outputPath = "../output.png"
const inputPath = "../assets/me.png"

func main() {
	img, err := utils.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	grayImg := utils.Grayscale(img)
	asciiArt := utils.GenerateASCIIArt(grayImg)

	fmt.Println(asciiArt)

	err = utils.RenderASCIIToImage(asciiArt, outputPath, 10)
	if err != nil {
		log.Fatalf("Failed to save ASCII art as PNG: %v", err)
	}

	fmt.Printf("Your ASCII art is ready, saved at %v.", outputPath)
}
