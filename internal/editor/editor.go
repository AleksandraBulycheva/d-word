package editor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AleksandraBulycheva/d-word/internal/file"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle     = lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230")).Padding(0, 1)
	statusBarStyle = lipgloss.NewStyle().Background(lipgloss.Color("236")).Foreground(lipgloss.Color("250")).Padding(0, 1)
	menuStyle      = lipgloss.NewStyle().Background(lipgloss.Color("235")).Foreground(lipgloss.Color("250")).Padding(0, 1)
	mainBoxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))
)

// Model represents state of editor
type Model struct {
	filename string
	fileSize int64
	textarea textarea.Model
	viewport viewport.Model
	width    int
	height   int
	isReady  bool
}

// New creates a new editor model
func New(filename string, content string) Model {
	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()
	ta.SetWidth(100)
	ta.SetHeight(20)

	var fileSize int64
	info, err := os.Stat(filename)
	if err == nil {
		fileSize = info.Size()
	}

	return Model{
		filename: filename,
		fileSize: fileSize,
		textarea: ta,
	}
}

// Init initializes the editor model
func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, taCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlS:
			err := file.WriteFile(m.filename, []byte(m.textarea.Value()))
			if err != nil {
				log.Fatal(err)
			}
			return m, tea.Quit
		case tea.KeyCtrlQ, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.isReady {
			m.viewport = viewport.New(msg.Width, msg.Height-10) // Adjust height for bars
			m.isReady = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 10 // Adjust height for bars
		}

		// Update textarea width within the box
		m.textarea.SetWidth(m.width - 4) // Account for box padding/border
	}

	m.viewport.SetContent(m.textarea.View())

	return m, tea.Batch(cmds...)
}

// View renders UI
func (m Model) View() string {
	if !m.isReady {
		return "Initializing..."
	}

	topBar := m.topBarView()
	menu := m.menuView()
	bottomBar := m.bottomBarView()

	editorBox := mainBoxStyle.Width(m.width - 2).Height(m.height - 8).Render(m.viewport.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		topBar,
		menu,
		editorBox,
		bottomBar,
	)
}

func (m *Model) topBarView() string {
	lineInfo := fmt.Sprintf("Line: %d/%d", 0, m.textarea.LineCount())
	sizeInfo := fmt.Sprintf("%.1fGB", float64(m.fileSize)/1.0e9)

	title := titleStyle.Render("d-wordedit v1.0")
	fileInfo := titleStyle.Render(fmt.Sprintf("%s (%s)", m.filename, sizeInfo))
	lineStatus := titleStyle.Render(lineInfo)

	w := lipgloss.Width

	totalWidth := w(title) + w(fileInfo) + w(lineStatus)
	gap := strings.Repeat(" ", max(0, m.width-totalWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, title, fileInfo, gap, lineStatus)
}

func (m *Model) menuView() string {
	menuItems := "[L]oad  [S]ave  [F]ind  [R]eplace  [E]dit  [T]ools  [H]elp  [Q]uit"
	return menuStyle.Width(m.width).Render(menuItems)
}

func (m *Model) bottomBarView() string {
	uniqueLines := "12,345 unique"
	memUsage := "Memory: 45.2MB/2.0GB"

	status := fmt.Sprintf("STATUS: %d lines | %s | %s", m.textarea.LineCount(), uniqueLines, memUsage)
	return statusBarStyle.Width(m.width).Render(status)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
