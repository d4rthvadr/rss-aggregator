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