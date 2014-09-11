package reddy

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"code.google.com/p/freetype-go/freetype"
	"github.com/jzelinskie/reddit"
	"github.com/quirkey/magick"
)

func NewPicturePoacher(subreddit string, session *reddit.LoginSession) *PicturePoacher {
	return &PicturePoacher{Subreddit: subreddit, Session: session}
}

type PicturePoacher struct {
	Subreddit            string
	Session              *reddit.LoginSession
	SubredditInfo        *reddit.Subreddit
	AcceptedImageDomains []string
}

func (pp *PicturePoacher) GetSubredditInfo() {
	var err error
	pp.SubredditInfo, err = pp.Session.AboutSubreddit(pp.Subreddit)
	if err != nil {
		fmt.Println(err)
	}
}

func (pp *PicturePoacher) writeImageScore() {

}

func (pp *PicturePoacher) FetchSubmissionImages() []ImageData {
	var resizeAll = "256x256"
	var cropAll = ""

	images := make([]ImageData, 0)

	submissions, _ := pp.Session.SubredditSubmissions(pp.Subreddit)

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
						if contentType[0] == "image" {
							data, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println("Unable to read body")
							} else {
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
							}
						}

					} else {
						fmt.Println("Failed to fetch image of type", imgType)
					}

				}
			}

		} else {
			fmt.Println("No image for ID", submissions[sub].ID)
		}
	}
	return images
}
