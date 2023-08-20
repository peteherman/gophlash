package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
	cardTextStyle = lipgloss.NewStyle().
			Width(30).
			Align(lipgloss.Center).
			PaddingLeft(5).
			PaddingRight(5).
			PaddingBottom(2)
	headerTextStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("200"))
	contextTextStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				PaddingLeft(1).
				Foreground(lipgloss.Color("100"))

	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("60"))
	focusTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("200"))
	softTitleStyle = lipgloss.NewStyle().
		PaddingLeft(2)
)

type model struct {
	library          Library
	deckIndex        int
	selectedDeckName string
	cardIndex        int
	showingFront     bool
	cursor           int
}

func initialModel(filepath string) model {
	library, err := ParseLibrary(filepath)
	if err != nil {
		fmt.Printf("Error reading library file: %v\n", err)
		os.Exit(1)
	}
	return model{
		library:          library,
		deckIndex:        0,
		selectedDeckName: "",
		cardIndex:        0,
		showingFront:     true,
		cursor:           0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.selectedDeckName == "" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "down", "ctrl+n":
				if m.deckIndex < len(m.library)-1 {
					m.deckIndex++
				}
			case "up", "ctrl+p":
				if m.deckIndex > 0 {
					m.deckIndex--
				}
			case "enter", " ":
				m.selectedDeckName = m.library[m.deckIndex].Name
			}
		}
		return m, nil
	} else {
		deck := m.library[m.deckIndex]
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "left", "down", "ctrl+p":
				if m.cardIndex > 0 {
					m.cardIndex--
					m.showingFront = true
				}
			case "right", "up", "ctrl+n":
				if m.cardIndex < len(deck.Cards)-1 {
					m.cardIndex++
					m.showingFront = true
				}
			case "enter", " ":
				m.showingFront = !m.showingFront
			}
		}
		return m, nil
	}
	return m, nil
}

func (m model) View() string {

	if m.selectedDeckName != "" {
		deck := m.library[m.deckIndex]
		header := softTitleStyle.Render(fmt.Sprintf("Deck: %v", deck.Name))
		header += "\n" + headerTextStyle.Render(fmt.Sprintf("Card: (%v/%v)", m.cardIndex+1, len(deck.Cards)))

		currentCard := deck.Cards[m.cardIndex]
		cardContent := ""
		if m.showingFront {
			cardContent += contextTextStyle.Render("Front")
			cardContent += "\n\n"
			cardContent += cardTextStyle.Render(currentCard.Front)
		} else {
			cardContent += contextTextStyle.Render("Back")
			cardContent += "\n\n"
			cardContent += cardTextStyle.Render(currentCard.Back)
		}
		help := helpTextStyle.Render("Use the keyboard to move cards,\nuse <space>/<enter> to flip cards,\npress q to quit.")
		return header + "\n" + cardStyle.Render(cardContent) + "\n" + help + "\n"
	}

	header := "Select the flashcard deck you'd like to review\n"
	body := ""
	for index, deck := range m.library {
		line := fmt.Sprintf("    %v", deck.Name)
		if m.deckIndex == index {
			line = focusTextStyle.Render(fmt.Sprintf(" >  %v", deck.Name))
		}
		body += line + "\n"
	}
	help := helpTextStyle.Render("Use the keyboard to move up/down,\nuse <space>/<enter> to select a deck,\npress q to quit.")	

	return header + body + help + "\n"
}

func main() {
	var libraryPath = flag.String("library", "", "Path to a .json file containing your library of gophlash cards")
	flag.Parse()

	if libraryPath != nil && *libraryPath == "" {
		if !DefaultLibraryExists() {
			fmt.Printf("No library file specified and the default library doesn't exist. Please specify where to find your library using the cmdline arg --library <path to library>.json\n")
			os.Exit(1)
		}
		fmt.Printf("Default library found. Using default library as no --library was specified\n")
		*libraryPath, _ = DefaultLibraryPath()
	}

	p := tea.NewProgram(initialModel(*libraryPath))
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v\n", err)
		os.Exit(1)
	}
}
