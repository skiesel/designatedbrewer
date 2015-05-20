package main

import (
	"github.com/andlabs/ui"
	"github.com/skiesel/designatedbrewer/screens"
	"os"
)

var window ui.Window

func init() {
	quitButton := ui.NewButton("Quit")
	quitButton.OnClicked(func() {
		window.Close()
		os.Exit(0)
	})

	tabs := ui.NewTab()
	tabs.Append("Home", quitButton)
	tabs.Append("Create", screens.GetCreateControl())
	tabs.Append("Load", screens.GetLoadControl())
	tabs.Append("Brew Day", screens.GetBrewControl())

	window = ui.NewWindow("DesignatedBrewer", 400, 500, tabs)
	window.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	window.Show()
}

func main() {
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
