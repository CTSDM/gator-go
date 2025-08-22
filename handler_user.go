package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return errors.New("No user found with that username.")
	} else if err != nil {
		return err
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Printf("The user %s\n has been logged in.\n", s.cfg.CurrentUserName)
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
	fmt.Printf("The user %s has been created.\n", user.Name)
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
