package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)
type Library []Deck

const DefaultLibraryFile string = "library.json"
const DefaultLibraryDirectory string = ".gophlash"

func DefaultLibraryDir() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDirectory, DefaultLibraryDirectory), nil
}

func DefaultLibraryPath() (string, error) {
	defaultDirectory, err := DefaultLibraryDir()
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

func SetupDefaultLibrary() error {
	defaultDir, err := DefaultLibraryDir()
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
	return library, nil
}
