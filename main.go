package main

import (
	"groupie-tracker/artists"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/artist", artists.Artists)
	mux.HandleFunc("/", artists.Home)
	mux.HandleFunc("/filtered", artists.Filtered)
	mux.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	log.Println("Запуск веб-сервера на http://localhost:8081/")
	err := http.ListenAndServe(":8081", mux)
	log.Fatal(err)
}
