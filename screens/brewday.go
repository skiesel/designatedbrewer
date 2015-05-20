package screens

import (
	"fmt"
	"github.com/andlabs/ui"
	"github.com/skiesel/designatedbrewer/sensors"
	"reflect"
	"time"
)

type temperatureRanges struct {
	meanValues []float64
	thresholds []float64
	active     []bool
}

var (
	brewSched brewSchedule
	brewTable ui.Table

	tempLabel0 ui.Label
	tempLabel1 ui.Label
	timeLabel  ui.Label

	timerStateChannel    = make(chan string, 1)
	timerDurationChannel = make(chan time.Duration, 1)

	temperatureStateChannel     = make(chan string, 1)
	temperatureThresholdChannel = make(chan temperatureRanges, 1)
)

func GetBrewControl() ui.Control {
	step := brewScheduleStep{}

	// row 1
	tempLabel0 = ui.NewLabel(formatTemperature(0))
	tempLabel1 = ui.NewLabel(formatTemperature(0))

	temperatures := ui.NewHorizontalStack(
		ui.NewLabel("Mash Temp"),
		tempLabel0,
		ui.NewLabel("Sparge Temp"),
		tempLabel1,
	)
	temperatures.SetStretchy(0)
	temperatures.SetStretchy(1)
	temperatures.SetStretchy(2)
	temperatures.SetStretchy(3)

	// row 2
	timeLabel = ui.NewLabel("00:00:00")
	timerControl := ui.NewHorizontalStack(ui.Space(), timeLabel, ui.Space())
	timerControl.SetStretchy(0)
	timerControl.SetStretchy(1)
	timerControl.SetStretchy(2)
	timerControl.SetPadded(true)

	// row 3
	startButton := ui.NewButton("start step")
	startButton.OnClicked(func() {
		timerDurationChannel <- time.Minute
		timerStateChannel <- "countdown"

		temperatureThresholdChannel <- temperatureRanges{active: []bool{true, true}}
		temperatureStateChannel <- "monitor"

		// timerStateChannel <- "countup"
	})
	autoAdvance := ui.NewCheckbox("auto-advance")

	stepControls := ui.NewHorizontalStack(
		startButton,
		autoAdvance,
	)
	stepControls.SetStretchy(0)
	stepControls.SetStretchy(1)
	stepControls.SetPadded(true)

	top := ui.NewVerticalStack(temperatures, timerControl, stepControls)
	top.SetStretchy(0)
	top.SetStretchy(1)
	top.SetStretchy(2)
	top.SetPadded(true)

	// bottom
	brewTable = ui.NewTable(reflect.TypeOf(step))

	stack := ui.NewVerticalStack(top, brewTable)

	stack.SetStretchy(0)
	stack.SetStretchy(1)

	go timerRoutine()
	go temperatureRoutine()

	return stack
}

func initBrewDay() {
	brewTable.Lock()
	d := brewTable.Data().(*[]brewScheduleStep)
	*d = brewSched.Steps
	brewTable.Unlock()
}

func intToPaddingString(n int64) string {
	if n > 0 {
		if n >= 10 {
			return fmt.Sprintf("%d", n)
		}
		return fmt.Sprintf("0%d", n)
	}
	return "00"
}

func durationToString(dur time.Duration) string {
	hours := int64(dur.Hours())
	minutes := int64(dur.Minutes()) % 60
	seconds := int64(dur.Seconds()) % 60

	hoursStr := intToPaddingString(hours)
	minutesStr := intToPaddingString(minutes)
	secondsStr := intToPaddingString(seconds)

	return fmt.Sprintf("%s:%s:%s", hoursStr, minutesStr, secondsStr)
}

func timerRoutine() {
	state := ""
	var timePoint time.Time

	for {
		if state == "" { //if not doing anything just block
			state = <-timerStateChannel
			timePoint = time.Now()
			if state == "countdown" {
				timerDuration := <-timerDurationChannel
				timePoint = timePoint.Add(timerDuration)
			}
		}
		select {
		case state = <-timerStateChannel:
			timePoint = time.Now()
			if state == "countdown" {
				timerDuration := <-timerDurationChannel
				timePoint = timePoint.Add(timerDuration)
			}
		default:
			if state == "" {
				continue
			}

			var dur time.Duration
			if state == "countup" {
				dur = time.Since(timePoint)
			} else if state == "countdown" {
				dur = timePoint.Sub(time.Now())
			}
			timeLabel.SetText(durationToString(dur))
			time.Sleep(time.Second)

		}
	}
}

func temperatureRoutine() {
	state := ""
	var ranges temperatureRanges
	for {
		if state == "" { //if not doing anything just block
			state = <-temperatureStateChannel
			ranges = <-temperatureThresholdChannel
		}

		select {
		case state = <-temperatureStateChannel:
			if state != "" {
				ranges = <-temperatureThresholdChannel
			}
		default:
			readings := sensors.GetThermometerReadings()
			for i, active := range ranges.active {
				if active {
					fmt.Printf("%d) %g\n", i, readings[i])
				}
			}
			time.Sleep(time.Second * 5)
		}

	}
}
