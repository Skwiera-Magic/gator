package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Skwiera-Magic/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_requests>", cmd.Name)
	}

	requestsInterval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid interval: %v", err)
	}

	log.Printf("Requesting feeds every %v", requestsInterval)

	ticker := time.NewTicker(requestsInterval)

	for ; ; <-ticker.C {
		requestFeeds(s)
	}
}

func requestFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("feed not found", err)
		return
	}
	log.Println("Feed found")
	collectFeed(s.db, feed)
}

func collectFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("could not mark feed %v fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("could not collect feed %v: %v", feed.Name, err)
		return
	}
	for _, data := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, data.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID: feed.ID,
			Title: data.Title,
			Description: sql.NullString{
				String: data.Description,
				Valid: true,
			},
			Url: data.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("could not create post: %v", err)
			continue
		}
	}
	log.Printf("%v posts found in feed %v", len(feedData.Channel.Item), feed.Name)
}