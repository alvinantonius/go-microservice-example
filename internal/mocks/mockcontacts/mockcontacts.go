package mockcontacts

import (
	"fmt"
	"github.com/alvinantonius/go-microservice-sample/internal/contacts"
)

type (
	// MockPkgContacts is mock object for PkgContacts struct
	MockPkgContacts struct{}

	// mockContact is mock object for contact struct
	mockContact struct {
		data contacts.ContactData
	}
)

// ReturnGet is the function that will be executed by MockPkgContacts.Get()
var ReturnGet func(int64) (contacts.Contact, error)

// ReturnCreate is the function that will be executed by MockPkgContacts.Create()
var ReturnCreate func(contacts.ContactData) (contacts.Contact, error)

// ReturnList is the function that will be executed by MockPkgContacts.List()
var ReturnList func(int64, int64) ([]contacts.ContactData, error)

// McUpdate is the function that will be executed by mocked contacts.Contact object
var McUpdate func(contacts.ContactData) error

// McDelete is the function that will be executed by mocked contacts.Contact object
var McDelete func() error

func init() {
	// initialize default ReturnGet function
	ReturnGet = func(contactID int64) (contacts.Contact, error) {
		cObj := mockContact{
			data: contacts.ContactData{
				ID: contactID,
			},
		}

		return &cObj, nil
	}

	// init default ReturnCreate function
	ReturnCreate = func(cData contacts.ContactData) (contacts.Contact, error) {
		cObj := mockContact{
			data: cData,
		}

		return &cObj, nil
	}

	// init default ReturnList function
	ReturnList = func(take, page int64) ([]contacts.ContactData, error) {
		start := (take * (page - 1)) + 1
		end := take * page

		var result []contacts.ContactData
		for i := start; i <= end; i++ {
			cData := contacts.ContactData{
				ID:    i,
				Name:  fmt.Sprintf("User%v", i),
				Email: fmt.Sprintf("email.user%v@example.com", i),
				Phone: fmt.Sprintf("+628123345678%v", i),
			}

			result = append(result, cData)
		}

		return result, nil
	}

	// init default McUpdate function
	McUpdate = func(input contacts.ContactData) error {
		return nil
	}

	// init default McDelete function
	McDelete = func() error {
		return nil
	}
}

// New will return MockPkgContacts for replacing PkgContacts object
func New() contacts.PkgContacts {
	return &MockPkgContacts{}
}

// Get is a mock function for PkgContacts.Get() function
func (mpc *MockPkgContacts) Get(contactID int64) (contacts.Contact, error) {
	return ReturnGet(contactID)
}

// Create is a mock function for PkgContacts.Create() function
func (mpc *MockPkgContacts) Create(input contacts.ContactData) (contacts.Contact, error) {
	return ReturnCreate(input)
}

// List is a mock function for PkgContacts.List() function
func (mpc *MockPkgContacts) List(take, page int64) ([]contacts.ContactData, error) {
	return ReturnList(take, page)
}

func (mc *mockContact) Update(input contacts.ContactData) error {
	return McUpdate(input)
}

func (mc *mockContact) Delete() error {
	return McDelete()
}

func (mc *mockContact) Data() contacts.ContactData {
	return mc.data
}
