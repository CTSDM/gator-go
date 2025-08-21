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

func (c *commands) nun(s *state, cmd command) error {
	err := c.handlers[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
