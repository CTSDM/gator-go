# Gator
Simple CLI to follow RSS and read the posts on the terminal.
It is part of a course project @ boot.dev

## Required
- Postgresql
- Go, at least version 1.23

## Installation
- Run `go install`.

## Setup
1. Create a postgresql database.
2. Create a config file `.gatorconfig.json` at your home directory. It should have a key `db_url` with value the url of the created database.
3. Run `gator-go` to open the CLI and verify it is installed. It should throw an error: `There should be at least one command.`

## Commands
`$ gator-go [command / parameters]`
- `$ gator-go register <username>`. Creates a user with name `username` and logs them in. A user can follow feeds and check them out.
- `$ gator-go login <username>`. Logs in the user `username`.
- `$ gator-go reset`. Delete all the users and feeds.
- `$ gator-go users`. Display all the users.
- `$ gator-go addfeed <name url>`. Add a feed name and url into the database and add the feed to the user's follow list. Only a logged in user can perform this action.
- `$ gator-go feeds`. Display all the feeds that are stored in the database.
- `$ gator-go agg <time>`. Fetches information periodically given the `<time>` provided. Time is in the form of i.e. `1h5m`.
- `$ gator-go follow <url>`. Add the feed to the user's follow list.
- `$ gator-go unfollow <url>`. Remove the feed from the user's follow list.
- `$ gator-go following`. Displays the follow list of the current user.
- `$ gator-go browse <limit>`. `limit` is an optional parameter to specify how many posts are shown; default value is 2. The most up to date posts are shown first.

## Usage example
```
$ gator-go register ctsdm
$ gator-go addfeed 'Hacker News' https://news.ycombinator.com/rss
$ gator-go follow https://news.ycombinator.com/rss
$ gator-go following
$ gator-go agg 1h # Press <Control>-C to cancel this process as it does not run on the background.
$ gator-go browse
```
