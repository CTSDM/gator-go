package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CTSDM/gator-go/internal/database"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("The 'agg' command takes one parameter, 'url'")
	}

	timeRefreshFeedInput := cmd.args[0]

	timeBetweenReqs, err := time.ParseDuration(timeRefreshFeedInput)
	if err != nil {
		return err
	}

	fmt.Println("--------------------------")
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}

	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	_, err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: currentTime, Valid: true},
		UpdatedAt:     currentTime,
		ID:            feed.ID,
	})
	if err != nil {
		return err
	}

	items := rss.Channel.Item
	printPostTitles(items, feed.Name)

	return nil
}

func printPostTitles(posts []RSSItem, feedName string) {
	fmt.Println("--------------------------")
	fmt.Printf("Posts from %s:\n", feedName)
	for _, post := range posts {
		fmt.Printf("- %s\n", post.Title)
	}
	fmt.Printf("--------------------------\n\n")
}
