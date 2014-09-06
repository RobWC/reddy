package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jzelinskie/reddit"
	"github.com/nfnt/resize"
)

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
	submissions, _ := session.SubredditSubmissions("nsfw")

	for sub := range submissions {
		if submissions[sub].Domain == "imgr.com" || submissions[sub].Domain == "i.imgur.com" {
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
						if contentType[0] == "image" {
							data, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println("Unable to read body")
							} else {
								fmt.Println(len(data))
								if strings.ToLower(contentType[1]) == "jpg" || strings.ToLower(contentType[1]) == "jpeg" {
									image, _, _ := image.Decode(bytes.NewReader(data))
									imgType = "jpg"
									newImage := resize.Resize(128, 128, image, resize.Lanczos3)
									images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
								} else if strings.ToLower(contentType[1]) == "png" {
									image, _, _ := image.Decode(bytes.NewReader(data))
									imgType = "png"
									newImage := resize.Resize(128, 128, image, resize.Lanczos3)
									images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
								} else if strings.ToLower(contentType[1]) == "gif" {
									image, _, _ := image.Decode(bytes.NewReader(data))
									imgType = "gif"
									newImage := resize.Resize(128, 128, image, resize.Lanczos3)
									images = append(images, ImageData{Image: newImage, Score: submissions[sub].Score})
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
	fmt.Println("Total files ", len(images))
}
