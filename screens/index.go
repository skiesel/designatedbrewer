package screens

import (
	"html/template"
	"net/http"
)

type page struct {
	Page string
	Data interface{}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		panic(err)
	}

	t.Execute(w, page{Page: "HOME"})
}
