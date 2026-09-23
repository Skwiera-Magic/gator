package main

import (
	"fmt"
	"context"
)

func reset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to reset database: %v", err)
	}
	fmt.Println("Database reset successfully")
	return nil
}