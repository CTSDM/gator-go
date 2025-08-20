package main

import (
	"errors"
	"fmt"

	"github.com/CTSDM/gator-go/internal/config"
)

type state struct {
	cfg config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) nun(s *state, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("The 'login' handler expects a single argument, 'username', but none was given!")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil
}
