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
)

type model struct {
	deck         Deck
	cardIndex    int
	showingFront bool
	cursor       int
}

func initialModel() model {
	deck := DeckFromFilepath("deck.json")
	return model{
		deck:         deck,
		cardIndex:    0,
		showingFront: true,
		cursor:       0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cardIndex < len(m.deck.Cards)-1 {
				m.cardIndex++
				m.showingFront = true
			}
		case "enter", " ":
			m.showingFront = !m.showingFront
		}
	}
	return m, nil
}

func (m model) View() string {
	header := "\n" + headerTextStyle.Render(fmt.Sprintf("Card: (%v/%v)", m.cardIndex+1, len(m.deck.Cards)))

	currentCard := m.deck.Cards[m.cardIndex]
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

func main() {
	var deckPath = flag.String("deck", "", "Path to a .json file containing the contents of your deck")
	flag.Parse()

	if deckPath != nil && *deckPath == "" {
		fmt.Printf("Please specify where to find your deck using the cmdline arg --deck <path to deck>.json\n")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v\n", err)
		os.Exit(1)
	}
}
