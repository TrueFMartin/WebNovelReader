package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"
	"log"
	"net/http"
	"strings"
)

// Holds information about chapter
type chapterContainer struct {
	title string
	body  []string
	url   string
}

// takes a url returns a goquery Document
func getDocFromURL(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func findTableOfContents() *goquery.Selection {
	doc := getDocFromURL("https://wanderinginn.com/table-of-contents/")
	return doc.Find(".entry-content")
}

// Creates array of chapters, sets chapter URLs and titles, returns chapter
func setBookURLs(selection *goquery.Selection) []chapterContainer {
	fmt.Printf("num Chapters: %v", len(selection.Nodes))
	allChapters := make([]chapterContainer, 0, len(selection.Nodes))

	//select all hyperlink nodes
	sel := selection.Find("a")
	sel.Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title and URL
		chapter := chapterContainer{}
		chapter.title = s.Text()
		url, ok := s.Attr("href")
		if ok == false {
			fmt.Println("no href")
			return
		}
		chapter.url = url
		allChapters = append(allChapters, chapter)
	})
	return allChapters
}

func getChapterBody(selection *goquery.Selection) []string {
	//select all <p> sections
	sel := selection.Find("p")
	chapter := make([]string, 0)

	sel.Each(func(i int, s *goquery.Selection) {
		paragraph := "<p>" + s.Text() + "</p>\n"
		chapter = append(chapter, paragraph)
	})
	return chapter
}

func main() {
	//Sets tile and author of ebook
	book := epub.NewEpub("The Wandering Inn")
	book.SetAuthor("Pirateaba")

	//Find ToC page, set as current selection
	selection := findTableOfContents()

	//Create chapter struct,
	//set chapter's URL from Table of Contents,
	//set chapter's Title
	chapterGroup := setBookURLs(selection)

	//Run through each chapter, adding body of chapter
	//then write chapter as section to ebook.epub
	for i, chapter := range chapterGroup {
		//set document from chapter's URL
		doc := getDocFromURL(chapter.url)
		//selection contains chapter's body
		selection = doc.Find(".entry-content")
		chapter.body = getChapterBody(selection)

		//Combine each paragraph in chapter's body[], leading with title
		var sectionStrBuilder strings.Builder //More efficient string concatenation
		sectionStrBuilder.WriteString("<h1>" + chapter.title + "</h1>")
		for _, paragraph := range chapter.body {
			sectionStrBuilder.WriteString(paragraph) //adds each paragraph to builder
		}

		//convert builder to string, add that section to ebook.epub
		_, err := book.AddSection(sectionStrBuilder.String(), chapter.title, "", "")
		if err != nil {
			log.Println("error section: " + chapter.title)
			log.Fatal(err)
		}
		//print progress of ebook writing
		percentComplete := fmt.Sprintf("%v/%v", i, len(chapterGroup)-1)
		fmt.Printf("\rProgress: %s", percentComplete)
	}

	//create and write .epub file

	err := book.Write("TheWanderingInn.epub")
	if err != nil {
		log.Fatal(err)
	}
}
