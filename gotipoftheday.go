package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const templatesDir = "templates"

type Day struct {
	DayNumber  int
	DayTopic   string
	DayContent template.HTML
}

func loadDay(w http.ResponseWriter, dayNumber int) {
	day := Day{
		DayNumber:  dayNumber,
		DayTopic:   "",
		DayContent: "",
	}

	var dayContentBuf bytes.Buffer

	tmpl, err := template.ParseFiles(templatesDir + fmt.Sprintf("/days/day%03d.html", dayNumber))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(&dayContentBuf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dayContent := dayContentBuf.String()
	splitDayContent := strings.SplitN(dayContent, "\n", 3)
	day.DayTopic = splitDayContent[0]
	day.DayContent = template.HTML(splitDayContent[2])

	tmpl, err = template.ParseFiles(templatesDir + "/days/base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dayHandler(w http.ResponseWriter, r *http.Request) {
	dayString := r.PathValue("id")
	dayNumber, err := strconv.Atoi(dayString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadDay(w, dayNumber)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templatesDir + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /day/{id}", dayHandler)
	mux.HandleFunc("GET /", indexHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
