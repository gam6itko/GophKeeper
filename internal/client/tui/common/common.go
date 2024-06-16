package common

type ListItem struct {
	title, desc string
	id          int
}

func NewListItem(title string, desc string, id int) ListItem {
	return ListItem{
		title: title,
		desc:  desc,
		id:    id,
	}
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }
func (i ListItem) ID() int             { return i.id }
