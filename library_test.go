package main

import (
	"os"
	"path/filepath"
	"testing"
)

var libraryContents = []byte("{\"decks\": {\"default\": {\"name\": \"default deck\",\"cards\": [{\"front\": \"Front of card 1\",\"back\": \"Back of card 1\"},{\"front\": \"Front of card 2\",\"back\": \"Back of card 2\"},{\"front\": \"Front of card 3\",\"back\": \"Back of card 3\"}]}}}")

func makeTempLibraryJSON(contents []byte) (string, error) {
	tmpDirectoryPath, err := os.MkdirTemp("", "testing-dir")
	if err != nil {
		return "", err
	}
	tmpFilepath := filepath.Join(tmpDirectoryPath, "tmp-test.json")
	if err = os.WriteFile(tmpFilepath, contents, 0444); err != nil {
		return "", err
	}
	return tmpFilepath, nil
}

func TestDefaultDirectory(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Error getting home directory for user: %v\n", err)
	}
	defaultDir, err := DefaultDirectory()
	if err != nil {
		t.Errorf("Received error when getting default dir: %v\n", err)
	}
	expectedDir := filepath.Join(homeDir, DefaultLibraryDirectory,
		DefaultLibraryFile)
	if expectedDir != defaultDir {
		t.Errorf("Expected default dir != received: %v, %v\n",
			expectedDir, defaultDir)
	}
}

func TestErrorWhenParseLibraryFromNonexistentFile(t *testing.T) {
	madeUpFilepath := "doesnt_exist.txt"
	_, err := ParseLibrary(madeUpFilepath)
	if err == nil {
		t.Errorf("Should've been an error for non-existent file\n")
		return
	}
}

func TestErrorWhenParseLibraryFromBadJSON(t *testing.T) {
	badJSON := []byte("This isn't good JSON\n")
	libraryFilepath, err := makeTempLibraryJSON(badJSON)
	if err != nil {
		t.Errorf("Error when creating temp library file\n")
		return
	}

	_, err = ParseLibrary(libraryFilepath)
	if err == nil {
		t.Errorf("Should've received error on malformed json library\n")
		return
	}
}

func TestLibraryFromGoodJSON(t *testing.T) {
	libraryFilepath, err := makeTempLibraryJSON(libraryContents)
	if err != nil {
		t.Errorf("Error when creating temp library file\n")
		return
	}

	library, err := ParseLibrary(libraryFilepath)
	if err != nil {
		t.Errorf("Received error when parsing library: %v\n", err)
	}
	if len(library.Decks) <= 0 {
		t.Errorf("Didn't parse any decks from json\n")
	}
}