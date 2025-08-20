package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CTSDM/gator-go/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reding file: %v", err)
	}
	runState := state{cfg: cfg}
	fmt.Printf("Read config: %v\n", runState.cfg)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("There should be at least one command.\n")
	}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	cmd := command{}
	cmd.name = args[1]
	cmd.args = []string{}
	if len(args) > 2 {
		cmd.args = args[2:]
	}

	err = cmds.nun(&runState, cmd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Reading modified config: %v\n", runState.cfg)
}
