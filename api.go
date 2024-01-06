package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type BookData struct {
	id     string
	title  string
	author string
	rating string
}

var fieldsOfInterest = map[string]bool{
	"field title": true,
}

func main() {
	resp, err := http.Get("https://www.goodreads.com/review/list/136622747-amy-armour?print=true&shelf=read&per_page=100")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	body := doc.FirstChild.LastChild

	table := findBooksTable(body)
	if table == nil {
		log.Fatal("Could not find book table")
	}

	for tableRow := table.FirstChild; tableRow != nil; tableRow = tableRow.NextSibling {
		// var data BookData
		for rowField := tableRow.FirstChild; rowField != nil; rowField = rowField.NextSibling {
			if isFieldOfInterest(rowField) {
				title := parseBookTitle(rowField)
				id := parseBookID(rowField)
				fmt.Printf("'%v'-'%v'\n", title, id)
			}
		}
	}
}

func parseBookTitle(rowField *html.Node) string {
	div := findFirstTag(rowField, "div")
	if div == nil {
		return ""
	}
	a := findFirstTag(div, "a")
	if a == nil {
		return ""
	}
	return strings.TrimSpace(a.FirstChild.Data)
}

func parseBookID(rowField *html.Node) string {
	div := findFirstTag(rowField, "div")
	if div == nil {
		return ""
	}
	a := findFirstTag(div, "a")
	if a == nil {
		return ""
	}
	for _, attr := range a.Attr {
		if attr.Key == "href" {
			re := regexp.MustCompile(`/book/show/(\d+)-`)
			m := re.FindStringSubmatch(attr.Val)
			if len(m) >= 2 {
				return m[1]
			}
		}
	}
	return ""
}

func findFirstTag(node *html.Node, tag string) *html.Node {
	for e := node.FirstChild; e != nil; e = e.NextSibling {
		if e.Type == html.ElementNode && e.Data == tag {
			return e
		}
	}
	return nil
}

func isFieldOfInterest(field *html.Node) bool {
	for _, attr := range field.Attr {
		key := attr.Key
		_, ofInterest := fieldsOfInterest[attr.Val]
		if key == "class" && ofInterest {
			return true
		}
	}
	return false
}

func findBooksTable(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}
	if node.Type == html.ElementNode &&
		node.Data == "tbody" &&
		len(node.Attr) != 0 &&
		node.Attr[0].Key == "id" &&
		node.Attr[0].Val == "booksBody" {

		return node
	}
	for nextNode := node.FirstChild; nextNode != nil; nextNode = nextNode.NextSibling {
		if res := findBooksTable(nextNode); res != nil {
			return res
		}
	}
	return nil
}
