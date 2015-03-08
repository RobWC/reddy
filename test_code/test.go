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
	"strconv"
	"strings"

	"code.google.com/p/freetype-go/freetype"
	"github.com/jzelinskie/geddit"
	"github.com/quirkey/magick"
)

var resizeAll = "256x256"
var cropAll = ""

func setText(newImage image.Image, newText string) image.Image {
	fontBytes, _ := ioutil.ReadFile("./FreeMonoBold.ttf")
	font, _ := freetype.ParseFont(fontBytes)
	rgba := image.NewRGBA(image.Rect(0, 0, newImage.Bounds().Max.X, newImage.Bounds().Max.Y))
	draw.Draw(rgba, rgba.Bounds(), newImage, image.ZP, draw.Src)
	textContext := freetype.NewContext()
	textContext.SetFontSize(32)
	textContext.SetFont(font)
	textContext.SetClip(newImage.Bounds())
	textContext.SetDst(rgba)
	textContext.SetSrc(image.White)
	pt := freetype.Pt(32, 32+int(textContext.PointToFix32(8)>>8))
	text := newText
	textContext.DrawString(text, pt)
	return rgba
}

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

	session, _ := geddit.NewLoginSession(
		"MrCrapperReddy",
		"XXX",
		"Reddy the Reddit Reader/1.0",
	)

	// Get reddit's default frontpage

	// Get our own personal frontpage
	submissions, _ := session.SubredditSubmissions("cats")
	//	subRedditInfo, _ := session.AboutSubreddit("aww")
	//fmt.Println(subRedditInfo)

	for sub := range submissions {
		if submissions[sub].Domain == "imgur.com" || submissions[sub].Domain == "i.imgur.com" {
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
									newImage, err := magick.NewFromBlob(data, "jpg")
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "jpg"
										newImage.Resize(resizeAll)
										newImage.Crop(cropAll)
										newImage.Strip()
										imageBlob, err := newImage.ToBlob("jpg")
										if err != nil {
											fmt.Println(err)
										} else {
											fontBytes, err := ioutil.ReadFile("./FreeMonoBold.ttf")
											font, err := freetype.ParseFont(fontBytes)
											finalImage, err := jpeg.Decode(bytes.NewReader(imageBlob))
											rgba := image.NewRGBA(image.Rect(0, 0, finalImage.Bounds().Max.X, finalImage.Bounds().Max.Y))
											draw.Draw(rgba, rgba.Bounds(), finalImage, image.ZP, draw.Src)
											textContext := freetype.NewContext()
											textContext.SetFontSize(32)
											textContext.SetFont(font)
											textContext.SetClip(finalImage.Bounds())
											textContext.SetDst(rgba)
											textContext.SetSrc(image.White)
											pt := freetype.Pt(32, 32+int(textContext.PointToFix32(8)>>8))
											text := strconv.FormatInt(int64(submissions[sub].Score), 10)
											textContext.DrawString(text, pt)

											if err != nil {

											} else {
												jpeg.Encode(bufio.NewWriter(bytes.NewBuffer([]byte{})), rgba, nil)
												images = append(images, ImageData{Image: rgba, Score: submissions[sub].Score})
											}
										}
									}
								} else if strings.ToLower(contentType[1]) == "png" {
									newImage, err := magick.NewFromBlob(data, "png")
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "png"
										newImage.Resize(resizeAll)
										newImage.Crop(cropAll)
										newImage.Strip()
										imageBlob, err := newImage.ToBlob("png")
										if err != nil {
											fmt.Println(err)
										} else {
											finalImage, err := png.Decode(bytes.NewReader(imageBlob))
											if err != nil {

											} else {
												images = append(images, ImageData{Image: finalImage, Score: submissions[sub].Score})
											}
										}
									}
								} else if strings.ToLower(contentType[1]) == "gif" {
									newImage, err := magick.NewFromBlob(data, "gif")
									if err != nil {
										fmt.Println(err)
									} else {
										imgType = "gif"
										newImage.Resize(resizeAll)
										newImage.Crop(cropAll)
										newImage.Strip()
										imageBlob, err := newImage.ToBlob("png")
										if err != nil {
											fmt.Println(err)
										} else {
											finalImage, err := gif.Decode(bytes.NewReader(imageBlob))
											if err != nil {

											} else {
												images = append(images, ImageData{Image: finalImage, Score: submissions[sub].Score})
											}
										}
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
	//maxHeigth := -1024
	finalImage := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	//draw.Draw(finalImage, finalImage.Bounds(), images[1].Image, image.ZP, draw.Src)
	//write out all images
	widthStart := -512
	heightStart := 0
	for item := range images {
		fmt.Println(widthStart, heightStart)
		imgBounds := images[item].Image.Bounds()
		if widthStart-imgBounds.Max.X <= maxWidth {
			widthStart = 0
			heightStart = heightStart - 128
		}
		draw.Draw(finalImage, finalImage.Bounds(), images[item].Image, image.Point{X: widthStart - widthStart - imgBounds.Max.X, Y: heightStart}, draw.Src)
		widthStart = widthStart - imgBounds.Max.X
		if widthStart <= maxWidth {
			widthStart = 0
			heightStart = heightStart - 128
		}
	}

	imageFile, _ := os.Create("FinalImage.png")
	defer imageFile.Close()
	imgWriter := bufio.NewWriter(imageFile)
	png.Encode(imgWriter, finalImage)
	fmt.Println("Total files ", len(images))
}
