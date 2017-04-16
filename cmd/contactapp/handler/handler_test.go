package handler

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alvinantonius/go-microservice-sample/internal/mocks/mockcontacts"

	"github.com/julienschmidt/httprouter"
)

func init() {
	pkgcontact = mockcontacts.New()
}

func TestListContact(t *testing.T) {
	method := "GET"

	testCase := []struct {
		Target    string
		Body      io.Reader
		ResStatus int
	}{
		{
			"http://www.example.com/v1/contacts?take=10",
			nil,
			200,
		},
		{
			"http://www.example.com/v1/contacts",
			nil,
			200,
		},
	}

	for index, tcase := range testCase {
		req := httptest.NewRequest(method, tcase.Target, tcase.Body)
		w := httptest.NewRecorder()
		p := httprouter.Params{}
		ListContact(w, req, p)

		resp := w.Result()

		if resp.StatusCode != tcase.ResStatus {
			t.Errorf("[TestListContact] tcase:%v res got %v | expect %v", index, resp.StatusCode, tcase.ResStatus)
		}
	}
}

func TestCreateContact(t *testing.T) {
	method := "POST"
	target := "http://www.example.com/v1/contacts"

	testCase := []struct {
		Body      io.Reader
		ResStatus int
	}{
		{
			strings.NewReader(`{"name":"User1", "email":"user1@email.com", "phone":"+628123456789"}`),
			201,
		},
		{
			strings.NewReader(`{"name":"Uer3 !@#" "email":"user1@email.com", "phone":"+628123456789"}`),
			400,
		},
	}

	for index, tcase := range testCase {
		req := httptest.NewRequest(method, target, tcase.Body)
		w := httptest.NewRecorder()
		p := httprouter.Params{}
		NewContact(w, req, p)

		resp := w.Result()

		if resp.StatusCode != tcase.ResStatus {
			t.Errorf("[TestCreateContact] tcase:%v res got %v | expect %v", index, resp.StatusCode, tcase.ResStatus)
		}
	}
}
