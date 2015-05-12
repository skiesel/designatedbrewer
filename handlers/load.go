package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type fileRequest struct {
	Filename string
}

func Load(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("saved")
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFiles("templates/load.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, page{ Page : "LOAD", Data : files })
}

func LoadFile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var fileName fileRequest
	err := decoder.Decode(&fileName)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadFile("saved/" + fileName.Filename)
	file := string(fileContents[:])
	fmt.Fprint(w, file)
}
