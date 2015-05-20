package screens

import (
	"bufio"
	"encoding/json"
	"github.com/andlabs/ui"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type brewSchedule struct {
	MashSteps   int64
	SpargeSteps int64
	BoilSteps   int64
	ChillSteps  int64
	Steps       []brewScheduleStep
}

type brewScheduleStep struct {
	Type        string
	Temperature string
	Duration    string
}

type createSchedule struct {
	MashSteps   int64
	SpargeSteps int64
	BoilSteps   int64
	ChillSteps  int64
	Steps       []createScheduleStep
}

type createScheduleStep struct {
	Type        string
	Temperature string
	Duration    string
	Delete      bool
}

type incrementControl struct {
	Buttons []ui.Button
	Labels  []ui.Label
	Value   int64
}

var (
	sched = createSchedule{
		MashSteps:   0,
		SpargeSteps: 0,
		BoilSteps:   0,
		ChillSteps:  0,
		Steps:       []createScheduleStep{},
	}

	scheduleTable ui.Table
)

func GetCreateControl() ui.Control {
	step := createScheduleStep{}

	scheduleTable = ui.NewTable(reflect.TypeOf(step))

	tabs := ui.NewTab()
	tabs.Append("Mash", getMashControl())
	tabs.Append("Sparge", getSpargeControl())
	tabs.Append("Boil", getBoilControl())
	tabs.Append("Chill", getChillControl())

	clear := ui.NewButton("delete selected steps")
	clear.OnClicked(func() {
		scheduleTable.Lock()
		filtered := sched.Steps[:0]
		size := 0
		sched.MashSteps = 0
		sched.SpargeSteps = 0
		sched.BoilSteps = 0
		sched.ChillSteps = 0
		for _, step := range sched.Steps {
			if !step.Delete {
				filtered = append(filtered, step)
				if strings.Contains(step.Type, "Mash") {
					sched.MashSteps += 1
					filtered[size].Type = "Mash Step " + strconv.FormatInt(sched.MashSteps, 10)
				} else if strings.Contains(step.Type, "Sparge") {
					sched.SpargeSteps += 1
					filtered[size].Type = "Sparge Step " + strconv.FormatInt(sched.SpargeSteps, 10)
				} else if strings.Contains(step.Type, "Boil") {
					sched.BoilSteps += 1
					filtered[size].Type = "Boil Step " + strconv.FormatInt(sched.BoilSteps, 10)
				} else if strings.Contains(step.Type, "Chill") {
					sched.ChillSteps += 1
					filtered[size].Type = "Chill Step " + strconv.FormatInt(sched.ChillSteps, 10)
				}
				size += 1
			}
		}
		sched.Steps = filtered
		d := scheduleTable.Data().(*[]createScheduleStep)
		*d = sched.Steps
		scheduleTable.Unlock()
	})

	save := ui.NewButton("save schedule")
	save.OnClicked(func() {
		now := time.Now()

		file, err := os.Create("saved/" + now.Format("2006-01-02_3:04PM"))
		if err != nil {
			panic(err)
		}

		defer file.Close()

		writer := bufio.NewWriter(file)
		encoder := json.NewEncoder(writer)
		encoder.Encode(sched)
		writer.Flush()
	})

	persistentButtons := ui.NewVerticalStack(clear, save)
	persistentButtons.SetPadded(true)

	stack := ui.NewVerticalStack(
		scheduleTable,
		tabs,
		persistentButtons,
	)

	stack.SetPadded(true)
	stack.SetStretchy(0)
	stack.SetStretchy(1)
	stack.SetStretchy(2)

	return stack
}

func getIncrementControls(min, max, startValue int64, label string, labelFunc func(int64) string, values []int64) *incrementControl {
	control := &incrementControl{
		Buttons: []ui.Button{},
		Labels: []ui.Label{
			ui.NewLabel(label),
			ui.NewLabel(labelFunc(startValue)),
		},
		Value: startValue,
	}

	for _, val := range values {
		button := ui.NewButton(strconv.FormatInt(val, 10))
		v := val
		button.OnClicked(func() {
			control.Value += v
			control.Labels[1].SetText(formatTemperature(control.Value))
		})

		control.Buttons = append(control.Buttons, button)
	}

	return control
}

func (ic incrementControl) getControls() []ui.Control {
	controls := []ui.Control{}
	for _, button := range ic.Buttons {
		controls = append(controls, button)
	}
	for _, label := range ic.Labels {
		controls = append(controls, label)
	}
	return controls
}

type IncrementCallback func(temperature, duration int64)

func getStepControls(temperature, duration bool, buttonLabel string, callable IncrementCallback) ui.Control {
	controls := []ui.Control{}

	tempIncrement := getIncrementControls(0, 212, 0, "temperature: ", formatTemperature, []int64{-10, -5, -1, 1, 5, 10})
	tempStack := ui.NewHorizontalStack(tempIncrement.getControls()...)
	tempStack.SetStretchy(6)
	tempStack.SetStretchy(7)
	tempStack.SetPadded(true)

	if temperature {
		controls = append(controls, tempStack)
	}

	timeIncrement := getIncrementControls(0, 120, 0, "duration: ", formatTime, []int64{-10, -5, -1, 1, 5, 10})
	timeStack := ui.NewHorizontalStack(timeIncrement.getControls()...)
	timeStack.SetStretchy(6)
	timeStack.SetStretchy(7)
	timeStack.SetPadded(true)

	if duration {
		controls = append(controls, timeStack)
	}

	button := ui.NewButton(buttonLabel)
	button.OnClicked(func() { callable(tempIncrement.Value, timeIncrement.Value) })
	controls = append(controls, button)

	stack := ui.NewVerticalStack(controls...)
	stack.SetPadded(true)

	return stack
}

func insertAt(index int64, step createScheduleStep, steps []createScheduleStep) []createScheduleStep {
	return append(steps[:index], append([]createScheduleStep{step}, steps[index:]...)...)
}

func getMashControl() ui.Control {
	var callable IncrementCallback = func(temperature, duration int64) {
		scheduleTable.Lock()
		sched.Steps = insertAt(sched.MashSteps, createScheduleStep{Type: "Mash Step " + strconv.FormatInt(sched.MashSteps+1, 10),
			Temperature: formatTemperature(temperature),
			Duration:    formatTime(duration),
		}, sched.Steps)
		sched.MashSteps += 1
		d := scheduleTable.Data().(*[]createScheduleStep)
		*d = sched.Steps
		scheduleTable.Unlock()
	}
	return getStepControls(true, true, "add mash step", callable)
}

func getSpargeControl() ui.Control {
	var callable IncrementCallback = func(temperature, duration int64) {
		scheduleTable.Lock()
		sched.Steps = insertAt(sched.MashSteps+sched.SpargeSteps,
			createScheduleStep{Type: "Sparge Step " + strconv.FormatInt(sched.SpargeSteps+1, 10),
				Temperature: formatTemperature(temperature),
				Duration:    formatTime(duration),
			}, sched.Steps)
		sched.SpargeSteps += 1
		d := scheduleTable.Data().(*[]createScheduleStep)
		*d = sched.Steps
		scheduleTable.Unlock()
	}
	return getStepControls(true, true, "add sparge step", callable)
}

func getBoilControl() ui.Control {
	var callable IncrementCallback = func(temperature, duration int64) {
		scheduleTable.Lock()
		sched.Steps = insertAt(sched.MashSteps+sched.SpargeSteps+sched.BoilSteps,
			createScheduleStep{Type: "Boil Step " + strconv.FormatInt(sched.BoilSteps+1, 10),
				Temperature: "--",
				Duration:    formatTime(duration),
			}, sched.Steps)
		sched.BoilSteps += 1
		d := scheduleTable.Data().(*[]createScheduleStep)
		*d = sched.Steps
		scheduleTable.Unlock()
	}
	return getStepControls(false, true, "add boil step", callable)
}

func getChillControl() ui.Control {
	var callable IncrementCallback = func(temperature, duration int64) {
		scheduleTable.Lock()
		sched.Steps = insertAt(sched.MashSteps+sched.SpargeSteps+sched.BoilSteps+sched.ChillSteps,
			createScheduleStep{Type: "Chill Step " + strconv.FormatInt(sched.ChillSteps+1, 10),
				Temperature: formatTemperature(temperature),
				Duration:    "--",
			}, sched.Steps)
		sched.ChillSteps += 1
		d := scheduleTable.Data().(*[]createScheduleStep)
		*d = sched.Steps
		scheduleTable.Unlock()
	}
	return getStepControls(true, false, "add chill step", callable)
}

func formatTemperature(t int64) string {
	return strconv.FormatInt(t, 10) + "°F"
}

func formatTime(t int64) string {
	return strconv.FormatInt(t, 10) + "min"
}

// func Create(w http.ResponseWriter, r *http.Request) {
// 	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
// 	if err == nil {
// 		t.Execute(w, page{Page: "CREATE"})
// 	}
// }

// func SaveSchedule(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var sched schedule
// 	err := decoder.Decode(&sched)
// 	if err != nil {
// 		panic(err)
// 	}

// 	now := time.Now()

// 	file, err := os.Create("saved/" + now.Format("2006-01-02_3:04PM"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer file.Close()

// 	writer := bufio.NewWriter(file)
// 	encoder := json.NewEncoder(writer)
// 	encoder.Encode(sched)
// 	writer.Flush()

// 	fmt.Fprint(w, "success")
// }
