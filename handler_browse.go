package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CTSDM/gator-go/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) > 0 {
		userLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Println("The specified limit of posts is not recognized. Defaulting to 2.")
		} else {
			limit = userLimit
		}
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)},
	)
	if err != nil {
		return err
	}

	printPosts(posts)
	return nil
}

func printPosts(posts []database.Post) {
	if len(posts) == 0 {
		fmt.Println("There are no posts to show.")
		return
	}

	fmt.Printf("--------------------------\n")
	fmt.Printf("Showing %v post(s):\n", len(posts))
	fmt.Printf("--------------------------\n")
	for i, post := range posts {
		fmt.Printf("Post number %v\n----->\n", i+1)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Published at: %s\n", post.PublishedAt.Time)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("----->\n\n")
	}
	fmt.Printf("--------------------------\n\n")
}
