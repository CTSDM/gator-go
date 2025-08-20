package config

import (
	"os"
	"testing"
)

func assertError(t testing.TB, err, want error) {
	t.Helper()

	if err != want {
		t.Errorf("err: %v, want: %v", err, want)
	}
}

func assertBool(t testing.TB, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func deleteTestFile(t testing.TB, pathToDelete string) error {
	t.Helper()

	err := os.Remove(pathToDelete)
	if err != nil {
		t.Error("couldnt delete the testfile", err)
	}
	return nil
}

func assertConfigStruct(t testing.TB, got, want Config, shouldMatch bool) {
	t.Helper()

	if shouldMatch {
		if got != want {
			t.Errorf("structures should match, got %v, want %v", got, want)
		}
	} else {
		if got == want {
			t.Errorf("structures should not match, got %v, want %v", got, want)
		}
	}
}
