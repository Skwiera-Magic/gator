package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Skwiera-Magic/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %v", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("could not get posts for user: %v", err)
	}

	fmt.Printf("Found %v posts for user %v:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%v from %v\n", post.PublishedAt.Time.Format("Mon 2 Jan"), post.FeedName)
		fmt.Printf("=== %v ===\n", post.Title)
		fmt.Printf("%v\n", post.Description.String)
		fmt.Printf("Link: %v\n", post.Url)
		fmt.Println("===========================")
	}
		
	return nil
}