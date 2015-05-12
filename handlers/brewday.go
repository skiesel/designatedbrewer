package handlers

import (
	"encoding/json"
	"github.com/skiesel/designatedbrewer/sensors"
	"html/template"
	"net/http"
	"os"
)

func BrewDay(w http.ResponseWriter, r *http.Request) {
	file := r.FormValue("schedule")
	fileReader, err := os.Open("saved/" + file)

	decoder := json.NewDecoder(fileReader)
	var sched schedule
	err = decoder.Decode(&sched)
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/brewday.html", "templates/header.html", "templates/footer.html")
	t.Execute(w, page{ Page : "BREW", Data : sched })
}

func GetTemperatureReadings(w http.ResponseWriter, r *http.Request) {
	readings := sensors.GetThermometerReadings()
	encoder := json.NewEncoder(w)
	encoder.Encode(readings)

}
