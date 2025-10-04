package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// mux := http.NewServeMux()
	log.Println("Starting HTMX Server ...")
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlstatment := `
	DROP TABLE IF EXISTS contacts;
	CREATE TABLE contacts(id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT, phone TEXT, email TEXT);
	`
	_, err = db.Exec(sqlstatment)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sqlite Setup Complete!")

	templates := template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
