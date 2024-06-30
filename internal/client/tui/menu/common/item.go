package common

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	_ list.Item = (*CmdItem)(nil)
)

type CmdItem struct {
	title, desc string
	cmd         tea.Cmd
}

func NewCmdItem(title string, desc string, cmd tea.Cmd) CmdItem {
	return CmdItem{
		title: title,
		desc:  desc,
		cmd:   cmd,
	}
}

func (i CmdItem) Title() string       { return i.title }
func (i CmdItem) Description() string { return i.desc }
func (i CmdItem) FilterValue() string { return i.title }
func (i CmdItem) Cmd() tea.Cmd        { return i.cmd }
