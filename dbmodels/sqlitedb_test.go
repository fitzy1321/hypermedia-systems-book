package dbmodels

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllContacts(t *testing.T) {
	db, err := GetAndSetupSqlite3DB("")
	assert.NoError(t, err)
	assert.NotNil(t, db, "DB is null")

	actual_contacts, err := db.GetAllContacts()
	assert.NoError(t, err)

	expected_contacts := []Contact{
		{Id: 1, FirstName: "Foo", LastName: "Bar", Phone: "Baz", Email: "Quux"},
	}
	assert.ElementsMatch(t, expected_contacts, actual_contacts)
}
