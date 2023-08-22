package library

import (
	"encoding/json"
	"fmt"
	"github.com/peteherman/gophlash/deck"
	"os"
	"path/filepath"
)

type Library struct {
	filepath string
	Decks    []deck.Deck
}

const DefaultLibraryFile string = "library.json"
const DefaultLibraryDirectory string = ".gophlash"

func DefaultLibrary() (Library, error) {
	if DefaultLibraryExists() {
		defaultLibraryPath, _ := DefaultLibraryPath()
		return LibraryFromFilepath(defaultLibraryPath)
	}
	return setupDefaultLibrary()
}

func LibraryFromFilepath(filepath string) (Library, error) {
	fileContents, err := os.ReadFile(filepath)
	if err != nil {
		return Library{}, err
	}
	var decks []deck.Deck
	err = json.Unmarshal(fileContents, &decks)
	if err != nil {
		return Library{}, err
	}

	library := Library{
		filepath: filepath,
		Decks:    decks,
	}

	return library, nil
}

func defaultLibraryDir() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, DefaultLibraryDirectory), nil
}

func DefaultLibraryPath() (string, error) {
	defaultDirectory, err := defaultLibraryDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(defaultDirectory, DefaultLibraryFile), nil
}

func DefaultLibraryExists() bool {
	libraryPath, err := DefaultLibraryPath()
	if err != nil {
		fmt.Printf("Error trying to determine your default library path: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(libraryPath); err == nil {
		return true
	}
	return false
}

func setupDefaultLibrary() (Library, error) {
	defaultDir, err := defaultLibraryDir()
	if err != nil {
		return Library{}, err
	}
	err = os.Mkdir(defaultDir, 0666)
	if err != nil && !os.IsExist(err) {
		return Library{}, err
	}
	library := Library{
		Decks: make([]deck.Deck, 0),
	}
	library.filepath, _ = DefaultLibraryPath()
	return library, nil
}
