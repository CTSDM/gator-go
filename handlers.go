package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CTSDM/gator-go/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("The 'login' handler expects a single argument, 'username', but none was given!")
	}

	// check if the users exists in the database
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return errors.New("No user found with that username.")
	} else if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Printf("The user %s\n has been logged in.", s.cfg.CurrentUserName)
	fmt.Println("--------------------------")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("The 'register' handler expect a single argument, 'username', but none was given!")
	}

	currentTime := time.Now()
	username := cmd.args[0]

	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Printf("The user %s\n has been created.", user.Name)
	fmt.Println("--------------------------")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DropUsers(context.Background())
	if err != nil {
		log.Println(err)
		return errors.New("Something went wrong while resetting the users table...")
	}

	fmt.Println("--------------------------")
	fmt.Println("The users table was reset successfully!")
	fmt.Println("--------------------------")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users in the database.")
		return nil
	}

	fmt.Println("--------------------------")
	fmt.Println("Current existing users:")
	for _, user := range users {
		toPrint := "* " + user.Name
		if user.Name == s.cfg.CurrentUserName {
			toPrint += " (current)"
		}
		fmt.Println(toPrint)
	}
	fmt.Println("--------------------------")

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	// fetches a RSS and prints the result
	url := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Println("Fetched RSS:")
	fmt.Println(*rss)
	fmt.Println("--------------------------")
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("The 'addfeed' command takes two args, 'name' and 'url'")
	}
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("The 'follow' command takes two args, 'name' and 'url'")
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

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
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

func handlerFollowing(s *state, cmd command) error {
	feed_follows, err := s.db.GetFeedfollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Printf("Current user: %s\n", s.cfg.CurrentUserName)
	fmt.Println("Followed feeds:")
	for _, feed_follow := range feed_follows {
		fmt.Printf("- %s\n", feed_follow.FeedName)
	}
	fmt.Println("--------------------------")

	return nil
}
