package deck

import (
	"encoding/json"
	"io/ioutil"
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

func DeckFromFilepath(f string) (Deck, error) {
	fileContents, err := ioutil.ReadFile(f)
	if err != nil {
		return Deck{}, err
	}

	deck, err := DeckFromFileContents(fileContents)
	if err != nil {
		return Deck{}, err
	}
	deck.filepath = f
	return deck, nil
}

func DeckFromFileContents(content []byte) (Deck, error) {
	deck := Deck{}
	err := json.Unmarshal(content, &deck)
	if err != nil {
		return Deck{}, err
	}
	return deck, nil
}
