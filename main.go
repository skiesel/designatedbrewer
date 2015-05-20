package main

import (
	"github.com/andlabs/ui"
	"github.com/skiesel/designatedbrewer/handlers"
	"github.com/skiesel/designatedbrewer/screens"
	"os"
	"net/http"
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

func webServer() {
	http.HandleFunc("/", handlers.Index)

	http.HandleFunc("/load", handlers.Load)
	http.HandleFunc("/load-file", handlers.LoadFile)

	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/save-schedule", handlers.SaveSchedule)

	http.HandleFunc("/brewday", handlers.BrewDay)
	http.HandleFunc("/get-temperature-readings", handlers.GetTemperatureReadings)

	http.HandleFunc("/send-alert", handlers.AlertMessage)

	http.HandleFunc("/push", handlers.Push)

	http.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir("sounds"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))
	http.Handle("/external/", http.StripPrefix("/external/", http.FileServer(http.Dir("external"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	go webServer()
	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
