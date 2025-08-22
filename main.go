package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/CTSDM/gator-go/internal/config"
	"github.com/CTSDM/gator-go/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reding file: %v", err)
	}

	// let's connect to the database
	db, err := sql.Open("postgres", cfg.DB_URL)

	dbQueries := database.New(db)
	runState := state{cfg: &cfg, db: dbQueries}

	args := os.Args
	if len(args) < 2 {
		log.Fatal("There should be at least one command.\n")
	}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("agg", handlerAggregator)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("users", handlerUsers)

	cmd := command{}
	cmd.name = args[1]
	cmd.args = []string{}
	if len(args) > 2 {
		cmd.args = args[2:]
	}

	err = cmds.run(&runState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
