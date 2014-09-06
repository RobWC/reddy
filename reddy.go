package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jzelinskie/reddit"
	"github.com/nfnt/resize"
)

type squareMask struct {
	p image.Point
}

func NewSquareMask(x int, y int) *squareMask {
	return &squareMask{p: image.Point{X: x, Y: y}}
}

func (sm *squareMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (sm *squareMask) Bounds() image.Rectangle {
	return image.Rect(sm.p.X, sm.p.Y, 128, 128)
}

func (sm *squareMask) At(x, y int) color.Color {
	xx, yy := float64(x-sm.p.X), float64(y-sm.p.Y)
	if xx*xx+yy*yy < float64(x*x+y*y) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

type ImageData struct {
	Image image.Image
	Score int
	Type  string
}

func main() {
	// Login to reddit
	images := make([]ImageData, 0)

	session, _ := reddit.NewLoginSession(
		"MrCrapperReddy",
		"eatTICKLE!@#",
		"Reddy the Reddit Reader/1.0",
	)

	// Get reddit's default frontpage

	// Get our own personal frontpage
	submissions, _ := session.SubredditSubmissions("aww")

	for sub := range submissions {
		if submissions[sub].Domain == "imgur.com" || submissions[sub].Domain == "i.imgur.com" {
			fmt.Println(submissions[sub].Score)
			_, err := url.Parse(submissions[sub].URL)
			imgType := ""
			if err != nil {
				fmt.Println("Unable to parse URL")
			} else {
				resp, err := http.Get(submissions[sub].URL)
				if err != nil {
					fmt.Println("Failed to fetch image")
				} else {
					if resp.StatusCode == 200 {
						//Check file type
						contentType := strings.Split(resp.Header["Content-Type"][0], "/")
						fmt.Println(contentType)
						if contentType[0] == "image" {
							data, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println("Unable to read body")
							} else {
								fmt.Println(len(data))
								if strings.ToLower(contentType[1]) == "jpg" || strings.ToLower(contentType[1]) == "jpeg" {
									image, err := jpeg.Decode(bytes.NewReader(data))
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "jpg"
										newImage := resize.Thumbnail(128, 128, image, resize.NearestNeighbor)
										images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
									}
								} else if strings.ToLower(contentType[1]) == "png" {
									image, err := png.Decode(bytes.NewReader(data))
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "png"
										newImage := resize.Thumbnail(128, 128, image, resize.NearestNeighbor)
										images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
									}
								} else if strings.ToLower(contentType[1]) == "gif" {
									image, err := gif.Decode(bytes.NewReader(data))
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "gif"
										newImage := resize.Thumbnail(128, 128, image, resize.NearestNeighbor)
										images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
									}
								}
								fmt.Println(imgType)
							}
						}

					} else {
						fmt.Println("Failed to fetch image")
					}

				}
			}

		} else {
			fmt.Println("No image for ", submissions[sub].ID)
		}
	}
	//Write out all images into a single image
	maxWidth := -1024
	finalImage := image.NewRGBA(image.Rect(0, 0, 1024, 512))
	//draw.Draw(finalImage, finalImage.Bounds(), images[1].Image, image.ZP, draw.Src)
	//write out all images
	widthStart := 0
	heightStart := 0
	for item := range images {
		fmt.Println(widthStart, heightStart)
		draw.Draw(finalImage, finalImage.Bounds(), images[item].Image, image.Point{X: widthStart, Y: heightStart}, draw.Src)
		widthStart = widthStart - 128
		if widthStart <= maxWidth {
			widthStart = 0
			heightStart = heightStart - 128
		}
	}
	//draw.Draw(finalImage, finalImage.Bounds(), images[2].Image, image.Point{X: -128, Y: 0}, draw.Src)
	//draw.Draw(finalImage, finalImage.Bounds(), images[3].Image, image.Point{X: -256, Y: 0}, draw.Src)

	//draw.DrawMask(finalImage, finalImage.Bounds(), images[1].Image, image.ZP, images[1].Image, image.ZP, draw.Over)
	//draw.DrawMask(finalImage, finalImage.Bounds(), images[2].Image, image.ZP, NewSquareMask(128, 256), image.ZP, draw.Over)

	imageFile, _ := os.Create("FinalImage.png")
	defer imageFile.Close()
	imgWriter := bufio.NewWriter(imageFile)
	png.Encode(imgWriter, finalImage)
	fmt.Println("Total files ", len(images))
}
