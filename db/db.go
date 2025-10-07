package db

import (
	"database/sql"
	"errors"
	"strings"
)

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

func ParseRowToContact(rows *sql.Rows) (*Contact, error) {
	var c Contact
	err := rows.Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Email)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (appdb *AppDB) GetAllContacts() ([]Contact, error) {
	rows, err := appdb.Query("SELECT id, first_name, last_name, phone, email FROM contacts;")
	if err != nil {
		return nil, err
	}
	var contacts []Contact
	for rows.Next() {
		c, err := ParseRowToContact(rows)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, errors.New("Contact could not be loaded")
		}
		contacts = append(contacts, *c)
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

func (appDB *AppDB) GetContactById(id int) (*Contact, error) {
	query := `SELECT id, first_name, last_name, phone, email FROM contacts WHERE id = ?;`
	var c Contact
	err := appDB.QueryRow(query, id).Scan(&c.Id, &c.FirstName, &c.LastName, &c.Phone, &c.Email)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (appDB *AppDB) UpdateContact(new *Contact) error {
	query := `
	UPDATE contacts
	SET first_name=?,
		last_name=?,
		phone=?,
		email=?
	WHERE id=?;`
	_, err := appDB.Exec(query, new.FirstName, new.LastName, new.Phone, new.Email, new.Id)
	return err
}

func (appDB *AppDB) DeleteContact(id int) error {
	query := `DELETE FROM CONTACTS WHERE id=?;`

	_, err := appDB.Exec(query, id)
	return err
}
