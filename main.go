package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

var output string
var err string

func main() {
	http.HandleFunc("/", GetHandler)
	http.HandleFunc("/ascii-art", PostHandler)
	go openBrowser("http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound, "")
		return
	}
	if r.Method != "GET" {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	tmpl, err := template.ParseFiles("Templates/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	err = tmpl.Execute(w, output)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "")
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "")
		return
	}

	input := r.FormValue("text")
	font := r.FormValue("fonts")
	input =strings.ReplaceAll(strings.ReplaceAll(input, "\\t", "    "), "\r", "\n")
	if input == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Please make sure you enter a text")
		return
	}

	font = "Fonts/" + font + ".txt"

	output, err = OutputArt(input, font)
	if err != "" {
		ErrorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, errM string) {
	var errorMessage string

	switch statusCode {
	case http.StatusNotFound:
		//404
		errorMessage = "Page not found"
	case http.StatusBadRequest:
		//400
		errorMessage = "Bad request"
		if errM != "" {
			errorMessage += ": " + errM
		}
	case http.StatusInternalServerError:
		//500
		errorMessage = "Internal server error"
	case http.StatusMethodNotAllowed:
		//405
		errorMessage = "Method not allowed"
	default:
		errorMessage = "Unexpected error"
	}

	data := struct {
		ErrorCode    int
		ErrorMessage string
	}{
		ErrorCode:    statusCode,
		ErrorMessage: errorMessage,
	}

	tmpl, err := template.ParseFiles("Templates/error.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	tmpl.Execute(w, data)
}

func openBrowser(URL string) {
	err := open.Run(URL)
	if err != nil {
		panic(err)
	}
}

func OutputArt(inputTXT, font string) (string, string) {
	outputstr := ""
	file, err := os.ReadFile(font)
	if err != nil {
		return "", "Please make sure you select a font"
	}
	var art []string
	m := ""
	for _, l := range file {
		if l == 10 {
			art = append(art, m)
			m = ""
		} else {
			m += string(l)
		}
	}

	inputTXTarray := strings.Split(inputTXT, "\n")

	for _, s := range inputTXTarray {
		for i := 1; i <= 8; i++ {
			for _, a := range s {
				if int(a) < 32 || int(a) > 126 {
					return "", "Invalid input"
				} else {
					outputstr += (art[(int(a)-32)*9+i])
				}
			}
			outputstr += "\n"
		}
		outputstr += "\n"
	}
	return outputstr, ""
}
