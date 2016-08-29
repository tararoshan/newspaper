package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

var (
	internRegex     = regexp.MustCompile(`(?i)[^\w](?:interns?|internship)[^\w]`)
	ansiEscapeRegex = regexp.MustCompile(`[[:cntrl:]]`)
	bayRegex        = regexp.MustCompile(`(?i)[^\w](?:san fran\w*|sf|bay area|mountain view|oakland|berkeley)[^\w]`)
	dallasRegex     = regexp.MustCompile(`(?i)[^\w](?:dallas|dfw|fort worth|richardson)[^\w]`)

	postTitleColor      = color.New(color.BgGreen).Add(color.Bold).SprintfFunc()
	commentTitleColor   = color.New(color.BgBlue).SprintfFunc()
	bayCommentColor     = color.New(color.FgMagenta).SprintfFunc()
	dallasCommentColor  = color.New(color.FgCyan).SprintfFunc()
	regularCommentColor = fmt.Sprintf
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func displayPaginatedString(paginatedString string, pager string) {
	cmd := exec.Command(pager)
	cmd.Stdin = strings.NewReader(paginatedString)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	panicIf(err)
}

func getLatestPost() (title string, url string) {
	whoishiring, err := goquery.NewDocument("https://news.ycombinator.com/submitted?id=whoishiring")
	panicIf(err)

	posts := whoishiring.Find("tr > td:nth-child(3)")
	for i := 0; i < posts.Length(); i++ {
		post := posts.Eq(i)
		title = post.Text()

		if strings.Contains(title, "Who is hiring?") {
			url, _ = post.Find("a").Attr("href")
			return title, url
		}
	}
	panic(errors.New("Error fetching/parsing whoishiring's posts."))
}

func getFormattedInternComments(url string) (formattedComments string, totalBay int, totalDallas int) {
	latestPost, err := goquery.NewDocument("https://news.ycombinator.com/" + url)
	panicIf(err)

	totalComments := 0
	comments := latestPost.Find("table .comtr > td > table > tbody > tr")
	for i := 0; i < comments.Length(); i++ {
		comment := comments.Eq(i)

		width, _ := comment.Find(".ind img").Attr("width")
		if width == "0" {
			text := ansiEscapeRegex.ReplaceAllStringFunc(comment.Text(), func(match string) string {
				if match == "\n" {
					return match
				}
				return ""
			})

			if internRegex.MatchString(text) {
				var formatComment func(string, ...interface{}) string

				switch {
				case bayRegex.MatchString(text):
					totalBay++
					formatComment = bayCommentColor
				case dallasRegex.MatchString(text):
					totalDallas++
					formatComment = dallasCommentColor
				default:
					formatComment = regularCommentColor
				}

				totalComments++
				formattedComments += "\n" // adding newlines to a template string with bgcolor makes it look bad
				formattedComments += commentTitleColor("Intern Comment No. %d", totalComments)
				formattedComments += "\n"
				formattedComments += formatComment("%s", text)
				formattedComments += "\n"
			}
		}
	}
	return
}

func main() {
	pager := os.Getenv("PAGER")

	title, url := getLatestPost()
	postTitle := postTitleColor("%s", title) + "\n"

	formattedComments, totalBay, totalDallas := getFormattedInternComments(url)
	postPrefix := bayCommentColor("Bay Area Comments (%d)", totalBay) + ", " + dallasCommentColor("Dallas Area Comments (%d)", totalDallas) + "\n"

	paginatedString := postTitle + postPrefix + formattedComments

	displayPaginatedString(paginatedString, pager)
}
