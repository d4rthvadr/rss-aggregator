package main

import (
	"time"

	"github.com/darthvadr/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func databaseToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
		ApiKey: dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	UserID    uuid.NullUUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedFollows struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.NullUUID  `json:"user_id"`
	FeedID    uuid.NullUUID  `json:"feed_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Post struct {
	ID          uuid.UUID     `json:"id"`
	URL         string        `json:"url"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	PublishedAt time.Time     `json:"published_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	FeedID      uuid.NullUUID  `json:"feed_id"`
}

func databaseToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		Title:     dbFeed.Title,
		URL:       dbFeed.Url,
		UserID:    dbFeed.UserID,
		CreatedAt: dbFeed.CreatedAt.Time,
		UpdatedAt: dbFeed.UpdatedAt.Time,
	}
}



func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feeds[i] = databaseToFeed(dbFeed)
	}
	return feeds
}

func databaseToFeedFollows(dbFeedFollows database.FeedFollow) FeedFollows {
	return FeedFollows{
		ID:        dbFeedFollows.ID,
		UserID:    dbFeedFollows.UserID,
		FeedID:    dbFeedFollows.FeedID,
		CreatedAt: dbFeedFollows.CreatedAt.Time,
		UpdatedAt: dbFeedFollows.UpdatedAt.Time,
	}
}

func databaseFeedFollowsToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollows {
	feedFollows := make([]FeedFollows, len(dbFeedFollows))
	for i, dbFeedFollows := range dbFeedFollows {
		feedFollows[i] = databaseToFeedFollows(dbFeedFollows)
	}
	return feedFollows
}

func databaseToPost(dbPost database.Post) Post {
	return Post{
		ID:          dbPost.ID,
		FeedID:      dbPost.FeedID,
		Description: dbPost.Description.String,
		URL:        dbPost.Url,
		Title:       dbPost.Title,
		PublishedAt: dbPost.PublishedAt,
		CreatedAt:   dbPost.CreatedAt.Time,
		UpdatedAt:   dbPost.UpdatedAt.Time,
	}
}


func databasePostsToPosts(dbPosts []database.Post) []Post {
	posts := make([]Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = databaseToPost(dbPost)
	}
	return posts
}