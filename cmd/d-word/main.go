package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AleksandraBulycheva/d-word/internal/editor"
	"github.com/AleksandraBulycheva/d-word/internal/file"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: d-word <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := file.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	m := editor.New(filename, string(content))

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
