package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

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
