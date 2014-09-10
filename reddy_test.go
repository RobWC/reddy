package reddy

import (
	"fmt"
	"testing"

	"github.com/jzelinskie/reddit"
)

func TestBasicSanity(t *testing.T) {
	session, _ := reddit.NewLoginSession(
		"MrCrapperReddy",
		"eatTICKLE!@#",
		"Reddy the Reddit Reader/1.0",
	)
	pp := NewPicturePoacher("aww", session)
	pp.GetSubredditInfo()
	images := pp.FetchSubmissionImages()
	fmt.Println(len(images))
	ig := NewImageGenerator(images)
	ig.CreateSquare(-1024, -1024)
	ig.SaveImage("testimage.png")
}
