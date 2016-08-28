package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var internRegex = regexp.MustCompile(`[^\w](?:interns?|internship)[^\w]`)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// Find the latest who's hiring post
	whoishiring, err := goquery.NewDocument("https://news.ycombinator.com/submitted?id=whoishiring")
	panicIf(err)

	var url string
	posts := whoishiring.Find("tr > td:nth-child(3)")
	for i := 0; i < posts.Length(); i++ {
		post := posts.Eq(i)
		title := post.Text()

		if strings.Contains(title, "Who is hiring?") {
			url, _ = post.Find("a").Attr("href")
			fmt.Println(title)
			break
		}
	}

	// Find and print all top level comments that match the intern regex
	latestPost, err := goquery.NewDocument("https://news.ycombinator.com/" + url)
	panicIf(err)

	totalComments := 0
	comments := latestPost.Find("table .comtr > td > table > tbody > tr")
	for i := 0; i < comments.Length(); i++ {
		comment := comments.Eq(i)

		width, _ := comment.Find(".ind img").Attr("width")
		if width == "0" {
			text := comment.Text()
			if internRegex.MatchString(text) {
				totalComments++
				fmt.Println("Intern Comment No.", totalComments)
				fmt.Println(text)
			}
		}
	}
}
