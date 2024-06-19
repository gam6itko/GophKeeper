package common

type WindowSize struct {
	Width  int
	Height int
}

type IWindowSizeAware interface {
	SetSize(width, height int)
}
