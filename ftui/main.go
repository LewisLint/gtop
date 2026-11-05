package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	dirs     []os.DirEntry
	selected int
}

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	dirs, err := os.ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	m := model{dirs: dirs}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func (m model) View() string {
	var out strings.Builder
	//	for _, entry := range m.dirs {
	//		name := entry.Name()
	//		if entry.IsDir() {
	//			name += "/"
	//		}
	//		fmt.Fprintln(&out, name) // dunno what fprintln does.
	//	}
	for i, entry := range m.dirs {
		prefix := "  "
		if i == m.selected {
			prefix = "> "
		}

		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		fmt.Fprintf(&out, "%s%s\n", prefix, name)
	}
	return out.String()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "q", "ctrl+c": // quit
			return m, tea.Quit
		case "k":
			if m.selected > 0 {
				m.selected--
			}
		case "j":
			if m.selected < len(m.dirs)-1 {
				m.selected++
			}
		}
	}
	return m, nil
}
