package main

import (
	"net/http"
	"golang.org/x/net/html"
	"fmt"
	"encoding/json"
)

type Page struct {
	Error bool
	ErrorMessage string
	Title string
	Posts [50]string
	PostLength int
}

func GetData(username string, page *Page) {
	page.Error = false
	content, err := http.Get("http://facebook.com/" + username  + "/posts")
	if err != nil {
		page.Error = true
		page.ErrorMessage = "Can't request page"
		return
	}

	tokenizer := html.NewTokenizer(content.Body)

	titleAdded := 0
	for {
		tag := tokenizer.Next()
		if tag == html.ErrorToken {
			return
		}

		token := tokenizer.Token()

		for _, a := range token.Attr {
			if a.Key == "data-testid" {
				// Deal with posts
				if a.Val == "post_message" {
					for i := 0; i < 5; i++ {
						tokenizer.Next()

						data := tokenizer.Token().Data
						if data == "div" {break}
						if data == "p" || data == "span" || data == "a" {continue}

						// Now check content
						page.Posts[page.PostLength] += data
					}

					page.PostLength++
				}
			}
		}

		if token.Data == "title" {
			tokenizer.Next()
			token = tokenizer.Token()

			// Only add the first title (facebook pages have
			// more than one title for some reason
			if titleAdded == 0 {
				page.Title = token.Data
			}

			titleAdded++
		}
	}

	if page.Title == "Page Not Found | Facebook" {
		page.Error = true
		page.ErrorMessage = "Username not found"
		return
	}
}

func httpGetUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"]
	if username == nil {
		fmt.Fprintf(w, "'No username'")
		return
	}

	page := Page{}
	GetData(username[0], &page)
	json, _ := json.Marshal(page)
	fmt.Fprintf(w, string(json))
}

func main() {
	http.HandleFunc("/", httpGetUsername)
	http.ListenAndServe(":8090", nil)
}
