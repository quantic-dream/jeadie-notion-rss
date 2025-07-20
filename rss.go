package main

import (
	"log"

	"github.com/mmcdole/gofeed"
)

func ParseRSS(url string) []*gofeed.Item {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		log.Printf("Error leyendo feed %s: %v", url, err)
		return nil
	}
	return feed.Items
}
