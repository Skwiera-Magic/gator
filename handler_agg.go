package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Skwiera-Magic/gator/internal/database"
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
		log.Println("could not mark feed %v fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Println("could not collect feed %v: %v", feed.Name, err)
		return
	}
	for _, data := range feedData.Channel.Item {
		fmt.Println("Found post: %v", data.Title)
	}
	log.Println("%v posts found in feed %v", len(feedData.Channel.Item), feed.Name)
}