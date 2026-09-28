package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Skwiera-Magic/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not retreive feed: %v", err)
	}

	feedRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %v", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollows(feedRow.UserName, feedRow.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not get feed follows: %v", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows for this user")
	}

	fmt.Printf("Feed follows for user %v:\n", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf("* %v\n", feedFollow.FeedName)
	}

	return nil
}

func printFeedFollows(username, feedname string) {
	fmt.Printf("* User: %v\n", username)
	fmt.Printf("* Feed: %v\n", feedname)
}