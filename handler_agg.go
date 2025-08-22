package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	// fetches a RSS and prints the result
	url := "https://www.wagslane.dev/index.xml"
	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Println("Fetched RSS:")
	fmt.Println(*rss)
	fmt.Println("--------------------------")
	return nil
}
