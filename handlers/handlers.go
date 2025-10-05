package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"text/template"

	hmsDB "github.com/fitzy1321/hypermedia-systems-book/db"
)

var baseTemplate = template.Must(template.ParseFiles("templates/_base.html"))

var tmpl = map[string]*template.Template{
	// "index":    template.Must(template.Must(baseTemplates.Clone()).ParseFiles("templates/index.html")),
	"contacts":       template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/contacts.html")),
	"new":            template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/new.html")),
	"contactdetails": template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/contactdetails.html")),
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

func Contacts(db *hmsDB.AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var contacts []hmsDB.Contact
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

func PostNewContact(appDB *hmsDB.AppDB) http.HandlerFunc {
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

		err = appDB.SaveContact(hmsDB.Contact{Id: 0, FirstName: first_name, LastName: last_name, Phone: phone, Email: email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/contacts?"+url.QueryEscape("flash=Created New User"), http.StatusTemporaryRedirect)
	}
}

func ContactDetails(appDB *hmsDB.AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid User Id", http.StatusBadRequest)
			return
		}

		contact, err := appDB.GetContactById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if contact == nil {
			http.NotFound(w, r)
			return
		}

		err = tmpl["contactdetails"].ExecuteTemplate(w, "base", map[string]any{
			"Title":   "Contact Details View",
			"Contact": contact,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
