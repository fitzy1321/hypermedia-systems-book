package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var baseTemplates = template.Must(template.ParseFiles("templates/base.html"))

var tmpl = map[string]*template.Template{
	"index":    template.Must(template.Must(baseTemplates.Clone()).ParseFiles("templates/index.html")),
	"contacts": template.Must(template.Must(baseTemplates.Clone()).ParseFiles("templates/contacts.html")),
}

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "contacts", http.StatusTemporaryRedirect)
		// data := map[string]any{
		// 	"Title": "Golang HTMX Contacts App",
		// }
		// err := tmpl["index"].ExecuteTemplate(w, "base", data)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
	}
}

type Contact struct {
	Id        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

func Contacts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// rows, err := db.Query("SELECT id, first_name, last_name, phone, email FROM contacts;")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }

		// var contacts []Contact
		// for rows.Next() {
		// 	var c Contact
		// 	err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Email)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	}
		// 	contacts = append(contacts, c)
		// }

		data := map[string]any{
			"Title": "Contacts Page",
			// "Contacts": contacts,
		}

		err := tmpl["contacts"].ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func GetAndSetupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}
	sqlstatment := `
	DROP TABLE IF EXISTS contacts;
	CREATE TABLE contacts(id INTEGER PRIMARY KEY, first_name TEXT, last_name TEXT, phone TEXT, email TEXT);
	`
	_, err = db.Exec(sqlstatment)
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}
	return db, nil
}

func main() {
	log.Println("Starting Sqlite3 ...")
	db, err := GetAndSetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Sqlite Setup Complete!")

	// Route Handler, prefer to use method routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", Index())
	mux.HandleFunc("GET /contacts", Contacts(db))

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

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
