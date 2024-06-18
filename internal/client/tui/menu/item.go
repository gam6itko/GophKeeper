package menu

import tea "github.com/charmbracelet/bubbletea"

type Item struct {
	title, desc string
	cmd         tea.Cmd
}

func NewItem(title string, desc string, cmd tea.Cmd) Item {
	return Item{
		title: title,
		desc:  desc,
		cmd:   cmd,
	}
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }
func (i Item) Cmd() tea.Cmd        { return i.cmd }
