package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ShortUrl struct {
	ID          int
	Shortened   string
	Original    string
	CreatedData time.Time
}

var urls []ShortUrl

func Home(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))

	data := struct {
		URLs []ShortUrl
	}{
		urls,
	}

	//fmt.Println(urls)

	tmpl.Execute(rw, data)
}

func ShortenURL(rw http.ResponseWriter, r *http.Request) {
	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(rw, "Por favor, ingrese una URL para acortarla.", http.StatusBadRequest)
		return
	}

	id := len(urls) + 1
	shortenedURL := fmt.Sprintf("http://localhost:8080/%d", id)

	urls = append(urls, ShortUrl{id, shortenedURL, originalURL, time.Now()})

	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 || id > len(urls) {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	index := id - 1
	if index < 0 || index >= len(urls) {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	shortURL := urls[index]

	http.Redirect(w, r, shortURL.Original, http.StatusMovedPermanently)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", Home)
	router.HandleFunc("/shorten", ShortenURL).Methods("POST")
	router.HandleFunc("/{id:[0-9]+}", Redirect)

	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println(urls)
}
