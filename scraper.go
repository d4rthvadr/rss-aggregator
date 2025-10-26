package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)



func startScraping(db *database.Queries, concurrency int, interval time.Duration) {

	log.Println("Starting scraper with concurrency:", concurrency, "and interval:", interval)
	
	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {

		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}

		wg := sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, feed, &wg)
		}

		wg.Wait()
	}
}


func scrapeFeed( db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Scraping feed:", feed.ID, feed.Url)

	// Update last_fetched_at immediately to avoid multiple workers fetching the same feed
	// and allow last_fetched_at to be updated once in this function
	err := db.UpdateFeedLastFetchedAt(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error updating feed last fetched at:", err)
		return
	}

	if feed.Url == "" {
		log.Println("Feed URL is empty, skipping")
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching RSS feed:", err)
		return
	}

	log.Printf("Fetched %d items from feed %s\n", len(rssFeed.Channel.Items), feed.Url)

	for _, item := range rssFeed.Channel.Items {

		log.Printf("Item: %s - %s\n", item.Title, item.Link)
	

		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing pubDate '%s': %v", item.PubDate, err)
			continue
		}

		// Here you would typically save the item to the database
		_, err = db.CreatePost(context.Background(), 
			database.CreatePostParams{
				ID: uuid.New(),
				Title: item.Title,
				Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
				Url: item.Link,
				Userid: feed.UserID,
				PublishedAt: parsedTime,
				CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
				UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
				FeedID: uuid.NullUUID{UUID: feed.ID, Valid: true},
			})
		if err != nil {
			log.Printf("Error creating post: %v", err)
		}
	}


	log.Printf("Finished scraping feed %d - %s\n", feed.ID, feed.Url)
}
