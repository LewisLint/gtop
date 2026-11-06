package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	dirs     []os.DirEntry
	selected int
	path     string
	message  string
	height   int
	offset   int
}

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	dirs, err := os.ReadDir(os.Args[1])
	if err != nil {
		panic(err)
	}
	m := model{dirs: dirs, path: os.Args[1]}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func (m model) View() string {
	var out strings.Builder

	end := m.offset + m.height

	if end > len(m.dirs) {
		end = len(m.dirs)

	}
	for i := m.offset; i < end; i++ {
		entry := m.dirs[i]

		n := entry.Name()
		if entry.IsDir() {
			n += "/"
		}
		color := ""
		kind := ml(entry)

		switch kind {
		case "l":
			color = "\x1b[36m"
		case "d":
			color = "\x1b[35m"
		case "x":
			color = "\x1b[32m"
		case "-":
			color = "\x1b[3m"
		}
		pfx := "  "
		if i == m.selected {
			pfx = "> "
		}
		fmt.Fprintf(&out, "%s\x1b[1m%s\x1b[22m%s %s\x1b[0m\n", pfx, color, kind, n)
	}
	return out.String()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.height = size.Height - 1
		if m.height < 1 {
			m.height = 1
		}
	}
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
		case "enter", "l":
			if len(m.dirs) == 0 {
				break
			}
			entry := m.dirs[m.selected]
			if !entry.IsDir() {
				break
			}
			npath := filepath.Join(m.path, entry.Name())
			ndir, err := os.ReadDir(npath)
			if err != nil {
				m.message = err.Error()
				break
			}
			m.path = npath
			m.dirs = ndir
			m.selected = 0
			m.message = ""

		case "h":

			npath := filepath.Dir(m.path)
			ndir, err := os.ReadDir(npath)
			if err != nil {
				m.message = err.Error()
				break
			}

			m.path = npath
			m.dirs = ndir
			m.selected = 0
			m.message = ""
		}
	}
	if m.height > 0 {
		if m.selected < m.offset {
			m.offset = m.selected
		}
		if m.selected >= m.offset+m.height {
			m.offset = m.selected - m.height + 1
		}
	}
	return m, nil
}

func ml(entry os.DirEntry) string {
	info, e := entry.Info()
	if e != nil {
		return "?"
	}

	mode := info.Mode()
	switch {
	case mode&os.ModeSymlink != 0:
		return "l"
	case mode.IsDir():
		return "d"
	case mode.Perm()&0111 != 0:
		return "x"
	default:
		return "-"
	}
}
