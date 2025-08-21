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

	fmt.Println("The user has been set")
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

	fmt.Printf("The new user %q was created\n", username)
	fmt.Println(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DropUsers(context.Background())
	if err != nil {
		log.Println(err)
		return errors.New("Something went wrong while resetting the users table...")
	}

	fmt.Println("The users table was reset successfully!")
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

	for _, user := range users {
		toPrint := "* " + user.Name
		if user.Name == s.cfg.CurrentUserName {
			toPrint += " (current)"
		}
		fmt.Println(toPrint)
	}

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(*rss)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("The 'addfeed' takes two args, 'name' and 'url'")
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

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}
