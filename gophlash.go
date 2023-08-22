package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peteherman/gophlash/library"
	"os"
	"strings"
)

var (
	staticDeckNameStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("200"))
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
	cardTitleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("100"))
	cardBodyStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			PaddingLeft(5).
			PaddingRight(5).
			Foreground(lipgloss.Color("255"))
	helpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("60"))
	focusTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("200"))
)

const (
	ViewMode = iota
	CreateMode
	EditMode
)

type keymap = struct {
	next, prev, add, quit, save, delete key.Binding
}

type model struct {
	library        library.Library
	deckIndex      int
	viewingDeck    bool
	cardIndex      int
	showingFront   bool
	cursor         int
	mode           int
	inputs         []textarea.Model
	textEditKeymap keymap
}

func main() {
	libraryPath := libraryPathFromCmdlineArgs()
	mode := modeFromCmdlineArgs()
	programModel := initialModel(libraryPath, mode)
	p := tea.NewProgram(programModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v\n", err)
		os.Exit(1)
	}
}

func libraryPathFromCmdlineArgs() string {
	var libraryPath = flag.String("library", "", "Path to a .json file containing your"+
		"library of gophlash cards")
	flag.Parse()

	if libraryPath != nil && *libraryPath == "" {
		checkDefaultLibraryExists()
		fmt.Printf("Default library found. Using default library as no " +
			"--library was specified\n")
		return *libraryPath
	}
	return *libraryPath
}

func checkDefaultLibraryExists() {
	if !library.DefaultLibraryExists() {
		fmt.Printf("No library file specified and the default library " +
			"doesn't exist. Please specify where to find your library " +
			"using the cmdline arg --library <path to library>.json\n")
		os.Exit(1)
	}
}

func modeFromCmdlineArgs() int {
	checkCmdlineArgsContainsSubcommand()
	subcommand := strings.ToLower(flag.Args()[0])
	switch subcommand {
	case "view":
		return ViewMode
	case "create":
		return CreateMode
	case "edit":
		return EditMode
	default:
		printSubcommandHelpMenu()
		return -1
	}
}

func checkCmdlineArgsContainsSubcommand() {
	if len(flag.Args()) < 1 {
		fmt.Printf("Please specify a subcommand, try 'help' to view a menu\n")
		os.Exit(1)
	}
}

func printSubcommandHelpMenu() {
	fmt.Printf("Gophlash is a way to view, create, and edit decks of flashcards " +
		"all within your terminal!\nSpecify one of the subcommands below to " +
		"get started\n" +
		"\tview\t\tView a deck of flashcards\n" +
		"\tcreate\t\tCreate a brand new deck of cards\n" +
		"\tedit\t\tEdit a deck of cards in your library\n")
	os.Exit(0)
}

func initialModel(libraryFilepath string, initialMode int) model {
	library := readLibrary(libraryFilepath)
	model := model{
		library:      library,
		deckIndex:    0,
		viewingDeck:  false,
		cardIndex:    0,
		showingFront: true,
		cursor:       0,
		mode:         initialMode,
		inputs:       make([]textarea.Model, 0),
		textEditKeymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			add: key.NewBinding(
				key.WithKeys("ctrl+plus"),
				key.WithHelp("ctrl+plus", "Append card to deck"),
			),
			save: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "Save deck"),
			),
			delete: key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "delete card"),
			),
		},
	}
	model.addInputs()
	return model
}

func readLibrary(filepath string) library.Library {
	library, err := library.LibraryFromFilepath(filepath)
	if err != nil {
		fmt.Printf("Error reading library file: %v\n", err)
		os.Exit(1)
	}
	return library
}

func (m model) addInputs() {
	m.inputs[0] = newDeckNameInput()
	m.inputs[1] = newCardInput()
	m.inputs[2] = newCardInput()

	if m.mode == CreateMode {
		m.inputs[0].Focus()
	}
}

func newDeckNameInput() textarea.Model {
	t := textarea.New()
	t.MaxHeight = 1
	t.MaxWidth = 1
	t.SetHeight(1)
	t.Placeholder = "Deck Name"
	return t
}
func newCardInput() textarea.Model {
	return textarea.New()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ViewMode:
		return updateViewMode(m, msg)
	case CreateMode:
		return updateCreateMode(m, msg)
	case EditMode:
		return m, nil
	default:
		return m, nil
	}
}

func updateViewMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.viewingDeck {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "down", "ctrl+n":
				if m.deckIndex < len(m.library.Decks)-1 {
					m.deckIndex++
				}
			case "up", "ctrl+p":
				if m.deckIndex > 0 {
					m.deckIndex--
				}
			case "enter", " ":
				m.viewingDeck = true
			}
		}
		return m, nil
	}
	deck := m.library.Decks[m.deckIndex]
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
		case "esc":
			m.viewingDeck = false
		}
	}
	return m, nil
}

func updateCreateMode(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.textEditKeymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit
		case key.Matches(msg, m.textEditKeymap.next):
			m.inputs[m.cursor].Blur()
			m.cursor++
			if m.cursor > len(m.inputs)-1 {
				m.cursor = 0
			}
			cmd := m.inputs[m.cursor].Cursor()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.textEditKeymap.prev):
			m.inputs[m.cursor].Blur()
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.inputs) - 1
			}
			cmd := m.inputs[m.cursor].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.textEditKeymap.add):
			fmt.Printf("Add pressed\n")
		case key.Matches(msg, m.textEditKeymap.delete):
			fmt.Printf("Delete pressed\n")
		case key.Matches(msg, m.textEditKeymap.save):
			fmt.Printf("Save pressed\n")			
		}
	}
	// Update all textareas
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)		

}

func (m model) View() string {
	switch m.mode {
	case ViewMode:
		return viewModeView(m)
	case CreateMode:
		return ""
	case EditMode:
		return ""
	default:
		return ""
	}
}

func viewModeView(m model) string {
	if m.viewingDeck {
		return viewModeDeckListView(m)
	}

	header := "Select the flashcard deck you'd like to review\n"
	body := ""
	for index, deck := range m.library.Decks {
		line := fmt.Sprintf("    %v", deck.Name)
		if m.deckIndex == index {
			line = focusTextStyle.Render(fmt.Sprintf(" >  %v", deck.Name))
		}
		body += line + "\n"
	}
	help := helpTextStyle.Render("Use the keyboard to move up/down,\nuse " +
		"<space>/<enter> to select a deck,\npress q to quit.")

	return header + body + help + "\n"
}

func viewModeDeckListView(m model) string {
	deck := m.library.Decks[m.deckIndex]
	header := staticDeckNameStyle.Render(fmt.Sprintf("Deck: %v", deck.Name))
	header += "\n" + staticDeckNameStyle.Render(fmt.Sprintf("Card: (%v/%v)",
		m.cardIndex+1, len(deck.Cards)))
	currentCard := deck.Cards[m.cardIndex]
	cardContent := ""
	if m.showingFront {
		cardContent += cardTitleStyle.Render("Front")
		cardContent += "\n\n"
		cardContent += cardBodyStyle.Render(currentCard.Front)
	} else {
		cardContent += cardTitleStyle.Render("Back")
		cardContent += "\n\n"
		cardContent += cardBodyStyle.Render(currentCard.Back)
	}
	help := helpTextStyle.Render("Use the keyboard to move cards,\nuse " +
		"<space>/<enter> to flip cards,\npress q to quit.")
	return header + "\n" + cardStyle.Render(cardContent) + "\n" + help + "\n"
}
