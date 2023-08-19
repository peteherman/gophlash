package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Deck struct {
	Name     string `json:name`
	filepath string
	Cards    []Card `json:cards`
}

type Card struct {
	Front string `json:front`
	Back  string `json:back`
}

func DeckFromFilepath(f string) Deck {
	fileContents, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Printf("Received error when reading Deck file: %v\n", err)
		os.Exit(1)
	}
	deck := DeckFromFileContents(fileContents)
	deck.filepath = f
	return deck
}

func DeckFromFileContents(content []byte) Deck {
	deck := Deck{}
	err := json.Unmarshal(content, &deck)
	if err != nil {
		fmt.Printf("Error parsing deck json contents: %v\n", err)
		os.Exit(1)
	}
	return deck
}
