package main

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

func (c *commands) run(s *state, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	// here we call the function that it's returned from middleware
	// we are doing effectively this
	// middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error
	// we are calling the function func(*state, command) error
	if err != nil {
		return err
	}
	return nil
}
