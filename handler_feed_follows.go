package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CTSDM/gator-go/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("The 'follow' command takes the parameter 'url'")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("there is no feed related to the given url.")
			return err
		}
		return err
	}

	// Check if the relation user-feed exists in the feed_follows table
	_, err = s.db.GetFeedFollow(context.Background(), database.GetFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err == nil {
		fmt.Printf("The feed %q is already being followed by the current user %q\n", feed.Name, user.Name)
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	currentTime := time.Now()

	input := database.CreateFeedFollowsParams{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	createdFeedFollow, err := s.db.CreateFeedFollows(context.Background(), input)
	if err != nil {
		return err
	}
	fmt.Println("--------------------------")
	fmt.Printf("Current user: %s\n", createdFeedFollow.UserName)
	fmt.Printf("Followed feed: %s\n", createdFeedFollow.FeedName)
	fmt.Println("--------------------------")

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	if len(feed_follows) == 0 {
		fmt.Println("--------------------------")
		fmt.Printf("%s is not following any feeds.\n", user.Name)
		fmt.Println("--------------------------")
		return nil
	}

	fmt.Println("--------------------------")
	fmt.Printf("Current user: %s\n", user.Name)
	fmt.Println("Followed feeds:")
	for _, feed_follow := range feed_follows {
		fmt.Printf("- %s\n", feed_follow.FeedName)
	}
	fmt.Println("--------------------------")

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("The 'follow' command takes the parameter 'url'")
	}
	url := cmd.args[0]

	_, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, Url: url})
	if err == sql.ErrNoRows {
		fmt.Printf("Can't unfollow feed %q because it's not in %s's follow list.\n", url, user.Name)
		return nil
	}
	return err
}
