package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

type RSSFeed struct {

	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
	
	
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {

	httpClient := http.Client{
		Timeout:  10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("Error fetching RSS feed: ", err)
		return RSSFeed{}, err

	}
	defer resp.Body.Close()

	var rssFeed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&rssFeed); err != nil {
		log.Println("Error decoding RSS feed: ", err)
		return RSSFeed{}, err
	}

	return rssFeed, nil
}