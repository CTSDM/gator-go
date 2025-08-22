package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CTSDM/gator-go/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("The 'addfeed' command takes two args, 'name' and 'url'")
	}
	feedName := cmd.args[0]
	Url := cmd.args[1]
	currentTime := time.Now()

	feedParams := database.CreateFeedParams{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      feedName,
		Url:       Url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowsParams{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollows(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Println(feed)
	fmt.Println("--------------------------")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Println("Existing feeds:")
	for _, feed := range feeds {
		fmt.Println(feed)
	}
	fmt.Println("--------------------------")
	return nil
}
