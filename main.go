package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB Stuff
func GetAndSetupDB() (*AppDB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}
	sqlstatment := `
	CREATE TABLE IF NOT EXISTS contacts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		phone TEXT,
		email TEXT
	);`
	_, err = db.Exec(sqlstatment)
	if err != nil {
		if db != nil {
			db.Close()
		}
		return nil, err
	}
	// Check if we need to insert sample data
	var count int
	db.QueryRow("SELECT COUNT(*) FROM contacts").Scan(&count)
	if count > 0 {
		return &AppDB{db}, nil
	}

	contacts := []struct {
		firstName, lastName, phone, email string
	}{
		{"Foo", "Bar", "Baz", "Foz"},
	}

	insertSQL := `INSERT INTO contacts (first_name, last_name, phone, email) VALUES (?,?,?,?);`
	for _, c := range contacts {
		_, err = db.Exec(insertSQL, c.firstName, c.lastName, c.phone, c.email)
		if err != nil {
			return nil, err
		}
	}

	return &AppDB{db}, nil
}

type Contact struct {
	Id        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

type AppDB struct {
	*sql.DB
}

func (appdb *AppDB) GetAllContacts() ([]Contact, error) {
	rows, err := appdb.Query("SELECT id, first_name, last_name, phone, email FROM contacts;")
	if err != nil {
		return nil, err
	}
	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Email)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (appdb *AppDB) GlobalSearchContacts(val string) ([]Contact, error) {
	val = strings.TrimSpace(val)
	searchterm := "%" + strings.ToLower(val) + "%"
	query := `
	SELECT id, first_name, last_name, phone, email
	FROM contacts
	WHERE LOWER(first_name) LIKE ?
		OR LOWER(last_name) LIKE ?
		OR LOWER(phone) LIKE ?
		OR LOWER(email) LIKE ?
	LIMIT 50;`
	rows, err := appdb.Query(query, searchterm, searchterm, searchterm, searchterm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []Contact
	for rows.Next() {
		var c Contact
		err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Email)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (appDB *AppDB) SaveContact(c Contact) error {
	insertSQL := `
		INSERT INTO contacts (first_name, last_name, phone, email)
		VALUES (?, ?, ?, ?);
	`
	_, err := appDB.Exec(insertSQL, c.FirstName, c.LastName, c.Phone, c.Email)
	return err
}

// end db stuff

// template stuff
var baseTemplate = template.Must(template.ParseFiles("templates/_base.html"))

var tmpl = map[string]*template.Template{
	// "index":    template.Must(template.Must(baseTemplates.Clone()).ParseFiles("templates/index.html")),
	"contacts": template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/contacts.html")),
	"new":      template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/new.html")),
}

// Handlers
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

func Contacts(db *AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contacts []Contact
		var err error
		// get optional "q" param
		query := r.URL.Query()
		q := query.Get("q")
		flash := query.Get("flash")
		if q != "" {
			contacts, err = db.GlobalSearchContacts(q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			contacts, err = db.GetAllContacts()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		data := map[string]any{
			"Contacts": contacts,
		}
		if q != "" {
			data["SearchTerm"] = q
		}
		if flash != "" {
			data["Flash"] = flash
		}

		err = tmpl["contacts"].ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func NewContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl["new"].ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func PostNewContact(appDB *AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		first_name := r.FormValue("first_name")
		last_name := r.FormValue("last_name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")

		err = appDB.SaveContact(Contact{Id: 0, FirstName: first_name, LastName: last_name, Phone: phone, Email: email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/contacts?"+url.QueryEscape("flash=Created New User"), http.StatusTemporaryRedirect)
	}
}

// end handlers

func main() {
	log.Println("Starting Sqlite3 ...")
	db, err := GetAndSetupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Sqlite Setup Complete!")

	mux := http.NewServeMux()
	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", Index())
	mux.HandleFunc("GET /contacts", Contacts(db))
	mux.HandleFunc("GET /contacts/new", NewContact())
	mux.HandleFunc("POST /contacts/new", PostNewContact(db))

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
