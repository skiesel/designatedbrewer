package screens

import (
	// "encoding/json"
	"github.com/andlabs/ui"
	// "github.com/skiesel/designatedbrewer/sensors"
	"reflect"
)

var (
	brewSched brewSchedule
	brewTable ui.Table
)

func GetBrewControl() ui.Control {
	step := brewScheduleStep{}

	brewTable = ui.NewTable(reflect.TypeOf(step))

	// refreshFileList()

	// selectButton := ui.NewButton("Load")

	// selectButton.OnClicked(LoadFile)

	stack := ui.NewVerticalStack(brewTable)

	return stack
}

// func BrewDay(w http.ResponseWriter, r *http.Request) {
// 	file := r.FormValue("schedule")
// 	fileReader, err := os.Open("saved/" + file)

// 	decoder := json.NewDecoder(fileReader)
// 	var sched schedule
// 	err = decoder.Decode(&sched)
// 	if err != nil {
// 		panic(err)
// 	}

// 	t, err := template.ParseFiles("templates/brewday.html", "templates/header.html", "templates/footer.html")
// 	t.Execute(w, page{Page: "BREW", Data: sched})
// }

// func GetTemperatureReadings(w http.ResponseWriter, r *http.Request) {
// 	readings := sensors.GetThermometerReadings()
// 	encoder := json.NewEncoder(w)
// 	encoder.Encode(readings)

// }
