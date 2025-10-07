// ***************************************************************
// * __                         ___            __                *
// */__  _  |  _. ._   _     |_| | |\/| \/    (_   _. | o _|_  _ *
// *\_| (_) | (_| | | (_| o  | | | |  | /\ o  __) (_| | |  |_ (/_*
// *                   _| /                /        |            *
// ***************************************************************
package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	hmsDB "github.com/fitzy1321/hypermedia-systems-book/db"
	"github.com/fitzy1321/hypermedia-systems-book/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting Sqlite3 ...")
	db, err := hmsDB.GetAndSetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Sqlite Setup Complete!")

	// templates
	var baseTemplate = template.Must(template.ParseFiles("templates/_base.html"))

	mux := http.NewServeMux()
	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", handlers.Index())
	mux.HandleFunc(
		"GET /contacts",
		handlers.Contacts(
			db,
			template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/view_table.html")),
		),
	)
	mux.HandleFunc(
		"GET /contacts/new",
		handlers.NewContact(
			template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/new.html"))),
	)
	mux.HandleFunc("POST /contacts/new", handlers.PostNewContact(db))
	mux.HandleFunc(
		"GET /contacts/{id}",
		handlers.ContactDetails(
			db,
			template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/view_details.html")),
		),
	)
	mux.HandleFunc(
		"GET /contacts/{id}/edit",
		handlers.GetContactEdit(
			db,
			template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/edit_form.html")),
		),
	)
	mux.HandleFunc("POST /contacts/{id}/edit", handlers.PostContactEdit(db))
	mux.HandleFunc("POST /contacts/{id}/delete", handlers.PostDeleteContact(db))

	// Server configuration
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
