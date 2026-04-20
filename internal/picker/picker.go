package picker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormatChoices returns display strings for file paths (basename only).
func FormatChoices(files []string) []string {
	choices := make([]string, len(files))
	for i, f := range files {
		choices[i] = filepath.Base(f)
	}
	return choices
}

// item implements list.Item for Bubbletea list.
type item struct {
	path string
	name string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.name }

// Pick shows an interactive picker and returns the selected file path.
// Falls back to numbered list if terminal is not interactive.
func Pick(files []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no image files found")
	}

	if !isInteractive() {
		return pickNumbered(files)
	}

	return pickBubbletea(files)
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func pickNumbered(files []string) (string, error) {
	names := FormatChoices(files)
	fmt.Println("Select an image to optimize:")
	for i, name := range names {
		fmt.Printf("  %d) %s\n", i+1, name)
	}

	fmt.Print("\n> ")
	var choice int
	if _, err := fmt.Scan(&choice); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	if choice < 1 || choice > len(files) {
		return "", fmt.Errorf("choice %d out of range (1-%d)", choice, len(files))
	}

	return files[choice-1], nil
}

type model struct {
	list     list.Model
	selected string
	quit     bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				m.selected = i.path
			}
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quit {
		return ""
	}
	return "\n" + m.list.View()
}

func pickBubbletea(files []string) (string, error) {
	items := make([]list.Item, len(files))
	for i, f := range files {
		items[i] = item{path: f, name: filepath.Base(f)}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("2")).
		BorderLeftForeground(lipgloss.Color("2"))

	l := list.New(items, delegate, 40, 14)
	l.Title = "Select an image to optimize"
	l.SetShowStatusBar(false)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Padding(0, 1)

	m := model{list: l}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("picker: %w", err)
	}

	final := finalModel.(model)
	if final.quit || final.selected == "" {
		return "", fmt.Errorf("cancelled")
	}

	return final.selected, nil
}
