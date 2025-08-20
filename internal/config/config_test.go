package config

import (
	"testing"
)

func TestWrite(t *testing.T) {
	t.Run("Return nil when correctly writing the json into a file", func(t *testing.T) {
		testPath := "./someRandomPath.json"
		dbURL := "postgres://example"

		input := Config{
			DB_URL: dbURL,
		}
		want := Config{
			DB_URL: dbURL,
		}

		err := write(testPath, input)
		got, err := read(testPath)

		assertConfigStruct(t, got, want, true)
		assertError(t, err, nil)

		deleteTestFile(t, testPath)
	})

	t.Run("It correctly adds the username to the json file", func(t *testing.T) {
		testPath := "./someRandomPath.json"
		dbURL := "postgres://example"
		username := "chano"

		want := Config{
			DB_URL:          dbURL,
			CurrentUserName: username,
		}

		err := write(testPath, want)
		assertError(t, err, nil)

		got, err := read(testPath)
		assertError(t, err, nil)

		assertConfigStruct(t, got, want, true)
		deleteTestFile(t, testPath)
	})

	t.Run("It correctly adds the username to the json file", func(t *testing.T) {
		testPath := "./someRandomPath.json"
		dbURL := "postgres://example"
		username := "chano"

		want := Config{
			DB_URL: dbURL,
		}

		err := write(testPath, want)
		assertError(t, err, nil)
		want.CurrentUserName = username + "random"

		got, err := read(testPath)
		assertError(t, err, nil)

		assertConfigStruct(t, got, want, false)
		deleteTestFile(t, testPath)
	})
}

func TestRead(t *testing.T) {
	// just check that the config file is loaded properly
	// we create a testfile to check if we properly read
	testPath := "./testfiles/testConfig.json"
	want := Config{
		DB_URL: "postgres://example",
	}

	got, err := read(testPath)
	assertError(t, err, nil)
	assertConfigStruct(t, got, want, true)
}
