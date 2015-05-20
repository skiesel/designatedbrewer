package screens

import (
	"encoding/json"
	"github.com/andlabs/ui"
	"io/ioutil"
	"reflect"
)

type scheduleFile struct {
	Filename string
}

var (
	scheduleFiles = []scheduleFile{}
	fileTable     ui.Table

	previewFile  brewSchedule
	previewTable ui.Table
)

func GetLoadControl() ui.Control {
	f := scheduleFile{}
	fileTable = ui.NewTable(reflect.TypeOf(f))
	fileTable.OnSelected(loadPreview)

	refreshFileList()

	s := brewScheduleStep{}
	previewTable = ui.NewTable(reflect.TypeOf(s))

	selectButton := ui.NewButton("Load")

	selectButton.OnClicked(loadBrewDayFile)

	stack := ui.NewVerticalStack(fileTable, previewTable, selectButton)

	stack.SetStretchy(0)
	stack.SetStretchy(1)
	stack.SetStretchy(2)

	return stack
}

func refreshFileList() {
	scheduleFiles = []scheduleFile{}

	files, err := ioutil.ReadDir("saved")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		scheduleFiles = append(scheduleFiles, scheduleFile{Filename: file.Name()})
	}

	fileTable.Lock()
	d := fileTable.Data().(*[]scheduleFile)
	*d = scheduleFiles
	fileTable.Unlock()
}

func loadPreview() {
	index := fileTable.Selected()
	if index < 0 {
		return
	}

	fileContents, err := ioutil.ReadFile("saved/" + scheduleFiles[index].Filename)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(fileContents, &previewFile)

	previewTable.Lock()
	d := previewTable.Data().(*[]brewScheduleStep)
	*d = previewFile.Steps
	previewTable.Unlock()
}

func loadBrewDayFile() {
	index := fileTable.Selected()
	if index < 0 {
		return
	}

	fileContents, err := ioutil.ReadFile("saved/" + scheduleFiles[index].Filename)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(fileContents, &brewSched)
	initBrewDay()
}
