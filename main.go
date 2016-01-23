package main

import (
	"github.com/sanjaybv/point/cserver"
	"github.com/sanjaybv/point/pserver"

	"html/template"
	"log"
	"net/http"
)

func main() {

	pserver.InitPointService()
	cserver.InitChatService()

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", hPoint)

	log.Println("starting Point on port 8080")
	log.Panic(http.ListenAndServe(":8080", nil))
}

func hPoint(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/point.html")
	if err != nil {
		log.Fatal(err)
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
