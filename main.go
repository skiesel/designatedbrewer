package main

import (
	"github.com/skiesel/designatedbrewer/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.Index)

	http.HandleFunc("/load", handlers.Load)
	http.HandleFunc("/load-file", handlers.LoadFile)

	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/save-schedule", handlers.SaveSchedule)

	http.HandleFunc("/brewday", handlers.BrewDay)
	http.HandleFunc("/get-temperature-readings", handlers.GetTemperatureReadings)

	http.HandleFunc("/send-alert", handlers.AlertMessage)

	http.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir("sounds"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))

	http.ListenAndServe(":8080", nil)
}
