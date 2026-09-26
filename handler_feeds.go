package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.SelectFeed(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Print(feed.Name)
		fmt.Print(feed.Url)
		username, err := s.db.GetUsername(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Print(username)
	}

	return nil
}