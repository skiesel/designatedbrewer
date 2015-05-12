package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

type stepPair struct {
	Temperature json.Number
	Duration    json.Number
}

type schedule struct {
	MashSteps   []stepPair
	SpargeSteps []stepPair
	BoilSteps   []json.Number
	ChillSteps  []json.Number
}

func Create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err == nil {
		t.Execute(w, page{Page: "CREATE"})
	}
}

func SaveSchedule(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var sched schedule
	err := decoder.Decode(&sched)
	if err != nil {
		panic(err)
	}

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

	fmt.Fprint(w, "success")
}
