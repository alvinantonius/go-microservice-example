package handler

import (
	"net/http"

	"github.com/alvinantonius/go-microservice-sample/internal/contacts"

	"github.com/julienschmidt/httprouter"
)

var pkgcontact contacts.PkgContacts

// Init handler
func Init() {
	pkgcontact = contacts.New()
}

// NewContact is for creating/insert new contact
func NewContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pkgcontact.Get(1)
}

// ListContact is for get list of contact
func ListContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// GetContact is for get 1 contact data by id
func GetContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// UpdateContact is for updating contact data
func UpdateContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// DeleteContact is for deleting 1 contact based on contact id
func DeleteContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
