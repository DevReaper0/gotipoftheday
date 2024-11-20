package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const templatesDir = "templates"

type Day struct {
	DayNumber  int
	DayContent string
}

func loadDay(w http.ResponseWriter, dayNumber int) {
	day := Day{
		DayNumber:  dayNumber,
		DayContent: "",
	}

	var dayContentBuf bytes.Buffer
	fmt.Fprintf(&dayContentBuf, "%03d", dayNumber)

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

	day.DayContent = dayContentBuf.String()

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
	dayString := r.PathValue("day")
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
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /day/{id}", dayHandler)
	mux.HandleFunc("GET /", indexHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Print("Server closed")
		} else {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down server")
	shoutDownCtx, shutdownRelease := context.WithTimeout(ctx, time.Second*10)
	defer shutdownRelease()
	if err := server.Shutdown(shoutDownCtx); err != nil {
		log.Fatal(err)
	}
}
