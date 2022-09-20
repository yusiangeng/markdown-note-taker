package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg Config

func main() {
	a := app.New()

	win := a.NewWindow("Markdown Note Taker")

	edit, preview := cfg.makeUi()

	win.SetContent(container.NewHSplit(edit, preview))

	win.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *Config) makeUi() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}
