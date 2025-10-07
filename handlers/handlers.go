package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	hmsDB "github.com/fitzy1321/hypermedia-systems-book/db"
)

// Handlers
func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "contacts", http.StatusSeeOther) // 303 redirct to GET
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

func Contacts(db *hmsDB.AppDB, tmpl *template.Template) http.HandlerFunc {
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
				http.Error(w, err.Error(), http.StatusInternalServerError) // 500
				return
			}
		} else {
			contacts, err = db.GetAllContacts()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError) // 500
				return
			}
		}
		// It's okay to have empty contacts list
		data := map[string]any{
			"Contacts": contacts,
		}
		if q != "" {
			data["SearchTerm"] = q
		}
		if flash != "" {
			data["Flash"] = flash
		}

		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}

	}
}

func NewContact(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
	}
}

func PostNewContact(appDB *hmsDB.AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
			return
		}

		first_name := r.FormValue("first_name")
		last_name := r.FormValue("last_name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")

		err = appDB.SaveContact(hmsDB.Contact{Id: 0, FirstName: first_name, LastName: last_name, Phone: phone, Email: email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		}

		http.Redirect(w, r, "/contacts?"+url.QueryEscape("flash=Created New User"), http.StatusSeeOther) // 303 redirect to GET
	}
}

func ContactDetails(appDB *hmsDB.AppDB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact, err := getContactFromPathID(appDB, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		if contact == nil {
			http.NotFound(w, r)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base", map[string]any{
			"Title":   "Contact Details View",
			"Contact": contact,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
	}
}

func GetContactEdit(appDB *hmsDB.AppDB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact, err := getContactFromPathID(appDB, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		if contact == nil {
			http.NotFound(w, r)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base", map[string]any{
			"Title":   "Edit Contact",
			"Contact": contact,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
	}
}

func PostContactEdit(appDB *hmsDB.AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact, err := getContactFromPathID(appDB, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		if contact == nil {
			http.NotFound(w, r)
			return
		}

		// get values from form
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // 400
			return
		}
		first_name := r.FormValue("first_name")
		if first_name == "" {
			first_name = contact.FirstName
		}
		last_name := r.FormValue("last_name")
		if last_name == "" {
			last_name = contact.LastName
		}
		email := r.FormValue("email")
		if email == "" {
			email = contact.Email
		}
		phone := r.FormValue("phone")
		if phone == "" {
			phone = contact.Phone
		}

		err = appDB.UpdateContact(
			&hmsDB.Contact{
				Id:        contact.Id,
				FirstName: first_name,
				LastName:  last_name,
				Phone:     phone,
				Email:     email,
			})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/contacts/%d", contact.Id), http.StatusSeeOther) // 303 redirect to GET
	}
}

func getContactFromPathID(appDB *hmsDB.AppDB, r *http.Request) (*hmsDB.Contact, error) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("Invalid Contact ID in Path.")
	}

	contact, err := appDB.GetContactById(id)
	if err != nil {
		return nil, err
	}
	return contact, nil
}

func PostDeleteContact(appDB *hmsDB.AppDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contact, err := getContactFromPathID(appDB, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
			return
		}
		if contact == nil {
			http.NotFound(w, r)
			return
		}

		err = appDB.DeleteContact(contact.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther) // 303 redirect to GET
	}
}
