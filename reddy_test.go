package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/jzelinskie/geddit"
)

func TestBasicSanityPics(t *testing.T) {
	session, _ := geddit.NewLoginSession(
		"MrCrapperReddy",
		"XXXXXXX",
		"Reddy the Reddit Reader/1.0",
	)
	pp := NewPicturePoacher("pics", session)
	pp.GetSubredditInfo()
	images := pp.FetchSubmissionImages()
	fmt.Println(len(images))
	ig := NewImageGenerator(images)
	ig.CreateSquare(-768, -768)
	ig.SaveImage("testimage-pics.png")
	log.Println("Saving image file...pics")
}

func TestBasicSanityAww(t *testing.T) {
	session, _ := geddit.NewLoginSession(
		"MrCrapperReddy",
		"XXXXXXX",
		"Reddy the Reddit Reader/1.0",
	)
	pp := NewPicturePoacher("aww", session)
	pp.GetSubredditInfo()
	images := pp.FetchSubmissionImages()
	fmt.Println(len(images))
	ig := NewImageGenerator(images)
	ig.CreateSquare(-768, -768)
	ig.SaveImage("testimage-aww.png")
	log.Println("Saving image file...aww")
}
