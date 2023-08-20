package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const DefaultLibraryFile string = "library.json"
const DefaultLibraryDirectory string = ".gophlash"

func DefaultDirectory() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, DefaultLibraryDirectory, DefaultLibraryFile), nil
}

func SetupDefaultLibrary() error {
	defaultDir, err := DefaultDirectory()
	if err != nil {
		return err
	}
	err = os.Mkdir(defaultDir, 0666)
	if err != nil && !os.IsExist(err) {
		return err
	}
	fmt.Printf("Creating default library directory (%v)\n", defaultDir)
	return nil
}

type Library struct {
	filepath string
	Decks    map[string]Deck `json:decks`
}

func ParseLibrary(filepath string) (Library, error) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return Library{}, err
	}
	var library Library
	err = json.Unmarshal(fileContents, &library)
	if err != nil {
		return Library{}, err
	}
	library.filepath = filepath
	return library, nil
}
