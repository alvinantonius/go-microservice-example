package contacts

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/alvinantonius/go-microservice-sample/internal/cache"
	"github.com/alvinantonius/go-microservice-sample/internal/database"

	"github.com/jmoiron/sqlx"
)

type (
	// PkgContacts object is used to call any method to get single contact object
	// or any method that doesn't require contact object
	PkgContacts interface {
		Get(int64) (Contact, error)
		List(int64, int64) ([]ContactData, error)
		Create(ContactData) (Contact, error)
	}

	// this struct is the main object of this package
	pkgContacts struct{}

	// Contact is the abstraction of contact object
	// we use interface to make it mockable and testable
	// since this object is created on runtime
	Contact interface {
		// For update data
		Update(ContactData) error

		// for delete data
		Delete() error

		// for get data
		Data() ContactData
	}

	// this struct is the main object of the package
	// will contains all method to manipulate single contact data
	contact struct {
		data     ContactData
		cacheKey string
	}

	// ContactData is the structure of one contact data
	ContactData struct {
		ID    int64  `json:"id" db:"id"`
		Name  string `json:"name" db:"name"`
		Email string `json:"email" db:"email"`
		Phone string `json:"phone" db:"phone"`
	}
)

var stmt map[string]*sqlx.Stmt

func prepareQueries() error {
	// if it's already prepared, don't prepare again
	if len(stmt) == 0 {
		stmt = make(map[string]*sqlx.Stmt)

		var err error

		dbconn, err := database.Conn("main", "slave")
		if err != nil {
			log.Panic("[contacts] fail get database connection -> ", err)
			return err
		}

		// Get 1 contact data from ID
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

		// Get many list of contact data
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
func New() PkgContacts {
	prepareQueries()
	return &pkgContacts{}
}

// Get contact by contact id
func (pkgc *pkgContacts) Get(contactID int64) (Contact, error) {

	// check cache first
	cacheKey := getCacheKey(contactID)

	// get cache data
	cacheConn, _ := cache.Conn("main")
	cacheMap, err := cacheConn.HGetAll(cacheKey).Result()
	if err != nil {
		log.Println("[Get] error get cache from redis ->", err)
	}

	// init empty object
	cObj := contact{cacheKey: cacheKey}
	cData := ContactData{}

	// if cache is empty, then we need to do query
	if len(cacheMap) == 0 {
		// get data from DB
		err := stmt["get"].QueryRowx(contactID).StructScan(&cData)
		if err != nil {
			log.Println("[Get] error get data from query ->", err)
			return nil, err
		}

		// prepare cache data
		cacheData := make(map[string]string)
		cacheData["id"] = strconv.FormatInt(cData.ID, 10)
		cacheData["name"] = cData.Name
		cacheData["email"] = cData.Email
		cacheData["phone"] = cData.Phone

		// store cache data
		err = cacheConn.HMSet(cacheKey, cacheData).Err()
		if err != nil {
			log.Println("[Get] fail store cache ->", err)
		}
	} else {
		if val, ok := cacheMap["id"]; ok {
			cData.ID, _ = strconv.ParseInt(val, 10, 64)
		}

		if val, ok := cacheMap["name"]; ok {
			cData.Name = val
		}

		if val, ok := cacheMap["email"]; ok {
			cData.Email = val
		}

		if val, ok := cacheMap["phone"]; ok {
			cData.Phone = val
		}
	}

	// fill data
	cObj.data = cData

	return &cObj, nil
}

// Create new contact
func (pkgc *pkgContacts) Create(input ContactData) (Contact, error) {
	// return if invalid
	if !validateContact(input) {
		return nil, errors.New("invalid contact data")
	}

	dbconn, _ := database.Conn("main", "master")

	tx, _ := dbconn.Beginx()
	defer tx.Rollback()

	var insertID int64

	err := tx.QueryRowx(`
			INSERT INTO
			contacts
				name,
				email,
				phone
			VALUES (
				$1,
				$2,
				$3
			) returning id
		`, input.Name, input.Email, input.Phone).Scan(&insertID)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("[Create] fail to commit ->", err)
		return nil, err
	}

	input.ID = insertID
	cObj := contact{data: input, cacheKey: getCacheKey(insertID)}

	return &cObj, nil
}

// List wil return list of contact data
func (pkgc *pkgContacts) List(take, page int64) ([]ContactData, error) {

	// validate input
	if take <= 0 || page <= 0 {
		return []ContactData{}, errors.New("Invalid input")
	}

	// calculate offset
	offset := take * (page - 1)

	rows, err := stmt["list"].Queryx(take, offset)
	if err != nil {
		log.Println("[List] error on query ->", err)
		return []ContactData{}, err
	}

	cList := []ContactData{}
	for rows.Next() {
		cData := ContactData{}
		rows.StructScan(&cData)
		cList = append(cList, cData)
	}

	return cList, nil
}

// Update contact data
func (c *contact) Update(input ContactData) error {
	return nil
}

func (c *contact) Delete() error {
	return nil
}

func (c *contact) Data() ContactData {
	return c.data
}

func validateContact(input ContactData) bool {
	return true
}

func getCacheKey(contactID int64) string {
	return fmt.Sprintf("contact:%v", contactID)
}
