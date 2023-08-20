package main

import (
	"os"
	"path/filepath"
	"testing"
)

var deckContents = []byte("{\"name\": \"Test Deck\",\"cards\": [{\"front\": \"Front of card 1\",\"back\": \"Back of card1\"}]}")

func makeTempDeckJSON(contents []byte) (string, error) {
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

func TestDeckFileInAnotherDirectory(t *testing.T) {
	tmpFilepath, err := makeTempDeckJSON(deckContents)
	if err != nil {
		t.Errorf("Error when creating tmpdeck for test: %v\n", err)
	}
	deck, err := DeckFromFilepath(tmpFilepath)
	if err != nil {
		t.Errorf("Error when creating deck: %v\n", err)
	}
	if deck.Name != "Test Deck" {
		t.Errorf("Error in name parsing of deck\n")
	}
	if len(deck.Cards) <= 0 {
		t.Errorf("Didn't parse any cards from the deck\n")
	}
}

func TestErrorWhenDeckMalformed(t *testing.T) {
	badContents := []byte("This is definitely not a deck\n")
	tmpFilepath, err := makeTempDeckJSON(badContents)
	if err != nil {
		t.Errorf("Error when creating tmpdeck for test: %v\n", err)
	}
	_, err = DeckFromFilepath(tmpFilepath)
	if err == nil {
		t.Errorf("Should've received error when creating this deck")
	}
}
