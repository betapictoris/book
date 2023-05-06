package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/taylorskalyo/goreader/epub"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/knipferrc/teacup/statusbar"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

const useHighPerformanceRenderer = false

var text string

// Bubble represents the properties of the UI.
type Bubble struct {
	statusbar statusbar.Bubble
	viewport  viewport.Model
	height    int
	content   string
	ready     bool

	book *epub.Rootfile
}

// Init initializes the UI.
func (Bubble) Init() tea.Cmd {
	return nil
}

// NewStatusbar creates a new instance of the UI.
func NewStatusbar() statusbar.Bubble {
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
		},
		statusbar.ColorConfig{
			Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
			Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
		},
	)

	return sb
}

// Update handles all UI interactions.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.height = msg.Height

		footerHeight := lipgloss.Height(b.footerView())
		verticalMarginHeight := footerHeight

		if !b.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			b.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			b.viewport.YPosition = 0
			b.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			b.viewport.SetContent(b.content)
			b.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			b.viewport.YPosition = 1
		} else {
			b.viewport.Width = msg.Width
			b.viewport.Height = msg.Height - verticalMarginHeight
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	// Handle keyboard and mouse events in the viewport
	b.viewport, cmd = b.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

// View returns a string representation of the UI.
func (b Bubble) View() string {
	if !b.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s", b.viewport.View(), b.footerView())
}

func (b Bubble) footerView() string {
	b.statusbar.SetSize(b.viewport.Width)
	b.statusbar.SetContent("Book CLI", b.book.Title, "", fmt.Sprintf("%3.f%%", b.viewport.ScrollPercent()*100))
	return b.statusbar.View()
}

func main() {
	fileName := ""

	if len(os.Args) >= 2 {
		fileName = os.Args[1]
	} else {
		log.Fatal("Usage: book <path to file>")
	}

	log.Info("Reading file...")

	reader, err := epub.OpenReader(fileName)
	if err != nil {
		log.Fatal("Failed to read", "error", err)
	}
	defer reader.Close() // Close the file when we are done with it
	book := reader.Rootfiles[0]

	log.Info("File read!", "title", book.Title)
	log.Info("Parsing file...")

	for _, item := range book.Manifest.Items {
		if item.MediaType == "application/xhtml+xml" {
			r, err := item.Open()
			if err != nil {
				log.Fatal("Failed to create item reader", "error", err)
			}
			cont, err := io.ReadAll(r)
			if err != nil {
				log.Fatal("Failed to read item", "error", err)
			}

			converter := md.NewConverter("", true, nil)
			content, err := converter.ConvertString(string(cont))
			if err != nil {
				log.Fatal("Failed to parse book content", "error", err)
			}

			text += content
		}
	}

	out, err := glamour.Render(text, "dark")
	if err != nil {
		log.Fatal("Failed to create glamour reader", "error", err)
	}

	p := tea.NewProgram(
		Bubble{statusbar: NewStatusbar(), content: string(out), book: book},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support, so we can track the mouse wheel
	)

	_, err = p.Run()
	if err != nil {
		log.Fatal("Could not run program", "error", err)
	}
}
