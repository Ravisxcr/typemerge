package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LogView struct {
	entry *widget.Entry
}

func NewLogView() *LogView {
	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord
	entry.Disable()
	entry.SetPlaceHolder("Generation messages will appear here.")
	return &LogView{entry: entry}
}

func (v *LogView) Object() fyne.CanvasObject {
	return v.entry
}

func (v *LogView) Reset(message string) {
	v.entry.SetText(message)
}

func (v *LogView) ShowLines(lines []string) {
	v.entry.SetText(strings.Join(lines, "\n"))
}
