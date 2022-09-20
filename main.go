package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg Config

const windowTitle = "Markdown Note Taker"

func main() {
	a := app.New()

	win := a.NewWindow(windowTitle)

	edit, preview := cfg.makeUi()
	cfg.createMenuItems(win)

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

func (app *Config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	saveAsMenuItem := fyne.NewMenuItem("Save as...", app.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var mdFileFilter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *Config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile == nil {
			return
		}

		write, err := storage.Writer(app.CurrentFile)
		if err != nil {
			dialog.ShowError(err, win)
		}

		write.Write([]byte(app.EditWidget.Text))
		defer write.Close()
	}
}

func (app *Config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// handle cancel button
			if read == nil {
				return
			}

			defer read.Close()

			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = read.URI()
			win.SetTitle(fmt.Sprintf("%s - %s", windowTitle, read.URI().Name()))
			app.SaveMenuItem.Disabled = false
		}, win)

		openDialog.SetFilter(mdFileFilter)
		openDialog.Show()
	}
}

func (app *Config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// handle cancel button
			if write == nil {
				return
			}

			// check filename ends in ".md"
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "File has to have .md extension", win)
				return
			}

			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(fmt.Sprintf("%s - %s", windowTitle, write.URI().Name()))
			app.SaveMenuItem.Disabled = false
		}, win)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(mdFileFilter)
		saveDialog.Show()
	}
}
