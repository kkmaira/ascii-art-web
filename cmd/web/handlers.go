package main

import (
	"html/template"
	"log"
	z1 "mkassymk/ascii-art-web-export-file/internal/functions"
	"net/http"
	"os"
	"strconv"
)

var Result string


func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
}

func AsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art/" {
		HandleError(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err.Error())
	}
	if !r.Form.Has("input") || !r.Form.Has("banner") {
		HandleError(w, http.StatusBadRequest)
		return
	}

	textToChange := r.FormValue("input")
	bannerToChange := r.FormValue("banner")

	output, err := z1.MakeAscii(textToChange, bannerToChange)

	Result= output
	if err != nil {
		HandleError(w, http.StatusBadRequest)
		return
	}

	err = tmpl.Execute(w, Result)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
}

func ExportFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}
	os.WriteFile("export.txt", []byte(Result), 0644)
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote("export.txt"))
	http.ServeFile(w, r, "./export.txt")
	
}

func HandleError(w http.ResponseWriter, num int) {
	errorData := struct {
		ErrorNum     int
		ErrorMessage string
	}{
		ErrorNum:     num,
		ErrorMessage: http.StatusText(num),
	}
	w.WriteHeader(num)
	tmpl, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, errorData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
