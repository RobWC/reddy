package reddy

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type ImageGenerator struct {
	images     []ImageData
	finalImage image.Image
}

func NewImageGenerator(newImages []ImageData) *ImageGenerator {
	return &ImageGenerator{images: newImages, finalImage: nil}
}

func (ig *ImageGenerator) CreateSquare(width int, heigth int) {
	finalImage := image.NewRGBA(image.Rect(0, 0, width, heigth))
	//draw.Draw(finalImage, finalImage.Bounds(), images[1].Image, image.ZP, draw.Src)
	//write out all images
	widthStart := 0
	heightStart := 0
	for item := range ig.images {
		fmt.Println(widthStart, heightStart)
		imgBounds := ig.images[item].Image.Bounds()
		if widthStart-imgBounds.Max.X <= width {
			widthStart = 0
			heightStart = heightStart - 128
		}
		draw.Draw(finalImage, finalImage.Bounds(), ig.images[item].Image, image.Point{X: widthStart, Y: heightStart}, draw.Src)
		widthStart = widthStart - imgBounds.Max.X
		if widthStart <= width {
			widthStart = 0
			heightStart = heightStart - 128
		}
	}
	ig.finalImage = finalImage
}

func (ig *ImageGenerator) SaveImage(filename string) {
	imageFile, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	imgWriter := bufio.NewWriter(imageFile)
	png.Encode(imgWriter, ig.finalImage)
	imageFile.Close()
}
