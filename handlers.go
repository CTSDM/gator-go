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

	userId := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}
	username := cmd.args[0]

	args := database.CreateUserParams{
		ID:        userId,
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
