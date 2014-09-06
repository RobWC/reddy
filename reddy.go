package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/jzelinskie/reddit"
	"github.com/nfnt/resize"
)

func fetchImage(imageURL string, path string, wg *sync.WaitGroup) image.Image {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Failed to fetch image")
	} else {
		if resp.StatusCode == 200 {
			//Check file type
			contentType := strings.Split(resp.Header["Content-Type"][0], "/")
			fileName := strings.TrimLeft(path, "/")
			if contentType[0] == "image" {
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				} else {
					newFile, err := os.Create(fileName)
					if err != nil {
						fmt.Println(err)
					} else {
						size, err := newFile.Write(data)
						if err != nil {
							fmt.Println(err)
						} else {
							err := newFile.Close()
							if err != nil {
								fmt.Println(err)
							} else {
								if strings.ToLower(contentType[1]) == "jpg" || strings.ToLower(contentType[1]) == "jpeg" {
									image, _, _ := image.Decode(bytes.NewReader(data))
									newImage := resize.Resize(128, 128, image, resize.Lanczos3)
									newImageFile, _ := os.Create(fmt.Sprintf("%s%s", fileName, "_small.jpg"))
									defer newImageFile.Close()
									smallImageWriter := bufio.NewWriter(newImageFile)
									jpeg.Encode(smallImageWriter, newImage, nil)
									smallImageWriter.Flush()
								} else if strings.ToLower(contentType[1]) == "jpg" {

								}
								fmt.Printf("Wrote file %s  at size %d checksum %x\n", fileName, size, sha1.Sum(data))
							}
						}
					}
				}
			}

		} else {
			fmt.Println("Error fetcing image code: ", resp.StatusCode)
		}

	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	// Login to reddit
	session, _ := reddit.NewLoginSession(
		"MrCrapperReddy",
		"eatTICKLE!@#",
		"Reddy the Reddit Reader/1.0",
	)

	// Get reddit's default frontpage

	// Get our own personal frontpage
	submissions, _ := session.SubredditSubmissions("funny")

	for sub := range submissions {
		if submissions[sub].Domain == "imgr.com" || submissions[sub].Domain == "i.imgur.com" {
			fmt.Println(submissions[sub].Score)
			subUrl, err := url.Parse(submissions[sub].URL)
			if err != nil {
				fmt.Println("Unable to parse URL")
			} else {
				wg.Add(1)
				go fetchImage(submissions[sub].URL, subUrl.Path, &wg)
			}

		} else {
			fmt.Println("No image for ", submissions[sub].ID, "\n\n")
		}
	}
	wg.Wait()
}
