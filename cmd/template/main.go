package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/drgrib/alfred"

	"github.com/drgrib/alfred-bear/core"
	"github.com/drgrib/alfred-bear/db"
)

func main() {
	text := os.Args[1]
	title := os.Args[2]
	tags := os.Args[3]
	dateTag := len(os.Args) > 4

	litedb, err := db.NewBearDB()
	if err != nil {
		panic(err)
	}

	content := core.ParseQuery(text)

	autocompleted, err := core.AutocompleteTags(litedb, content)
	if err != nil {
		panic(err)
	}

	if !autocompleted {
		callback := []string{}
		if title != "" {
			callback = append(callback, "title="+url.PathEscape(title))
		} else {
			callback = append(callback, "title="+url.PathEscape(
				time.Now().Format("2006-01-02")))
		}

		tags := strings.Split(tags, ",")
		if dateTag {
			tags = append(tags, time.Now().Format("2006/01/02"))
		}
		if len(content.Tags) != 0 {
			bareTags := []string{}
			for _, t := range content.Tags {
				bareTags = append(bareTags, url.PathEscape(t[1:]))
			}
			tags = append(tags, bareTags...)
		}
		if len(tags) > 0 {
			callback = append(callback, "tags="+strings.Join(tags, ","))
		}

		if content.WordString == "" {
			clipString, err := clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
			if clipString != "" {
				callback = append(callback, "text="+url.PathEscape(clipString))
			}
		} else {
			str, err := ioutil.ReadFile(content.WordString)
			if err != nil {
				panic(err)
			}
			if len(str) > 0 {
				callback = append(callback, "text="+url.PathEscape(string(str)))
			}
		}

		callbackString := strings.Join(callback, "&")

		item := alfred.Item{
			Title: fmt.Sprintf("Create %#v", content.WordString),
			Arg:   callbackString,
			Valid: alfred.Bool(true),
		}
		if len(content.Tags) != 0 {
			item.Subtitle = strings.Join(content.Tags, " ")
		}
		alfred.Add(item)
	}

	alfred.Run()
}

