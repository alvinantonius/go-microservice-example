package contacts

import (
	"log"

	// "github.com/alvinantonius/go-microservice-sample/internal/cache"
	"github.com/alvinantonius/go-microservice-sample/internal/database"

	"github.com/jmoiron/sqlx"
)

type (
	// PkgContacts object is used to call any method to get single contact object
	// or any method that doesn't require contact object
	PkgContacts struct{}

	// Contact is the abstraction of contact object
	Contact interface {
		Update(ContactData) error
		Delete() error
	}

	// this struct is the main object of the package
	// will contains all method to manipulate single contact data
	contact struct {
		data ContactData
	}

	// ContactData is the structure of one contact data
	ContactData struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`

		cacheKey string
	}
)

var stmt map[string]*sqlx.Stmt

func prepareQueries() error {
	log.Println("call")
	if len(stmt) == 0 {
		log.Println("run")
		stmt = make(map[string]*sqlx.Stmt)

		var err error

		dbconn, err := database.Conn("main", "slave")
		if err != nil {
			log.Panic("[contacts] fail get database connection -> ", err)
			return err
		}

		stmt["get"], err = dbconn.Preparex(`
		SELECT
			id, name, email, phone
		FROM 
			contacts
		WHERE id = $1
	`)
		if err != nil {
			log.Println(err)
		}

		stmt["list"], err = dbconn.Preparex(`
		SELECT
			id, name, email, phone
		FROM
			contacts
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2
	`)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

// New will return contact struct as Contact interface
// if this run in test, it will return mocked contact struct
func New() *PkgContacts {
	prepareQueries()
	return &PkgContacts{}
}

// Get contact by contact id
func (pkgc *PkgContacts) Get(contactID int64) (Contact, error) {
	return &contact{}, nil
}

// Create new contact
func (pkgc *PkgContacts) Create(input ContactData) (Contact, error) {
	return &contact{}, nil
}

// List wil return list of contact data
func (pkgc *PkgContacts) List(take, page int64) ([]ContactData, error) {
	return []ContactData{}, nil
}

// Update contact data
func (c *contact) Update(input ContactData) error {
	return nil
}

func (c *contact) Delete() error {
	return nil
}
