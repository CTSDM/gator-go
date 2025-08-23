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
	err = savePosts(s, items, feed)
	return err
}

func savePosts(s *state, posts []RSSItem, feed database.Feed) error {
	fmt.Println("--------------------------")
	fmt.Printf("Saving posts from %s...\n", feed.Name)
	if len(posts) == 0 {
		fmt.Println("There are no posts to save...")
		fmt.Println("--------------------------")
		return nil
	}

	newPostsSaved := 0
	for _, post := range posts {
		pubDate, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			return err
		}

		currentTime := time.Now()
		newPostParams := database.CreatePostParams{
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
			Title:       post.Title,
			Url:         sql.NullString{String: post.Link, Valid: true},
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: pubDate, Valid: true},
			FeedID:      feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), newPostParams)
		if err == nil {
			newPostsSaved++
			continue
		} else if err != sql.ErrNoRows {
			return err
		}
	}

	if newPostsSaved == 0 {
		fmt.Printf("From %v fetched posts, no new posts were found.\n", len(posts))
	} else {
		fmt.Printf("From %v fetched posts, %v posts were new and were successfully saved.\n", len(posts), newPostsSaved)
	}
	fmt.Printf("--------------------------\n\n")
	return nil
}
