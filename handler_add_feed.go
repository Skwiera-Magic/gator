package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Skwiera-Magic/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		Name: name,
		Url: url,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follows: %v", err)
	}

	fmt.Println("Feed created:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("Feed followed:")
	printFeedFollows(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("======================")
	
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID: %s\n", feed.ID)
	fmt.Printf("* Created: %v\n", feed.CreatedAt)
	fmt.Printf("* Updated: %v\n", feed.UpdatedAt)
	fmt.Printf("* Name: %s\n", feed.Name)
	fmt.Printf("* URL: %s\n", feed.Url)
	fmt.Printf("* UserID: %s\n", feed.UserID)
}