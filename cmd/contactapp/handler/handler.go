package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alvinantonius/go-microservice-sample/internal/contacts"

	"github.com/julienschmidt/httprouter"
)

var pkgcontact contacts.PkgContacts

type (
	// Response is a struct that used to return JSON object for all request
	Response struct {
		Error interface{} `json:"errors,omitempty"`
		Links interface{} `json:"links,omitempty"`
		Data  interface{} `json:"data,omitempty"`
	}
)

// Init handler
func Init() {
	pkgcontact = contacts.New()
}

// NewContact is for creating/insert new contact
func NewContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// get json input data
	decoder := json.NewDecoder(r.Body)
	var input contacts.ContactData
	err := decoder.Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = pkgcontact.Create(input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return

}

// ListContact is for get list of contact
func ListContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	take, _ := strconv.ParseInt(r.FormValue("take"), 10, 64)
	page, _ := strconv.ParseInt(r.FormValue("page"), 10, 64)

	// set default take and max take
	if take == 0 || take >= 100 {
		take = 5
	}

	// set default page
	if page == 0 {
		page = 1
	}

	data, err := pkgcontact.List(take, page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepare result header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// write result
	res := Response{Data: data}
	jsonByte, _ := json.Marshal(res)
	w.Write(jsonByte)

	return
}

// GetContact is for get 1 contact data by id
func GetContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	contactID, _ := strconv.ParseInt(p.ByName("contact_id"), 10, 64)
	if contactID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cObj, err := pkgcontact.Get(contactID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepare result header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// write result
	res := Response{Data: cObj.Data()}
	jsonByte, _ := json.Marshal(res)
	w.Write(jsonByte)

	return
}

// UpdateContact is for updating contact data
func UpdateContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get param contact id
	contactID, _ := strconv.ParseInt(p.ByName("contact_id"), 10, 64)
	if contactID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get json input data
	decoder := json.NewDecoder(r.Body)
	var input contacts.ContactData
	err := decoder.Decode(&input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cObj, err := pkgcontact.Get(contactID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if contactID not found
	if cObj == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = cObj.Update(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepare result header
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// write result
	res := Response{Data: cObj.Data()}
	jsonByte, _ := json.Marshal(res)
	w.Write(jsonByte)

	return
}

// DeleteContact is for deleting 1 contact based on contact id
func DeleteContact(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get param contact id
	contactID, _ := strconv.ParseInt(p.ByName("contact_id"), 10, 64)
	if contactID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cObj, err := pkgcontact.Get(contactID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = cObj.Delete()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
