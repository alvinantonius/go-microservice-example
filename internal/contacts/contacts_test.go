package contacts

import (
	"errors"
	"log"
	"reflect"
	"testing"

	"github.com/alvinantonius/go-microservice-sample/internal/cache"
	"github.com/alvinantonius/go-microservice-sample/internal/database"

	"github.com/alicebob/miniredis"
	"github.com/jmoiron/sqlx"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// sqlmock obj
var mock sqlmock.Sqlmock
var pkgCon PkgContacts
var prepared map[string]*sqlmock.ExpectedPrepare

// this init function will only called when run unit test
func init() {
	// create new sqlmock obj
	db, mockinit, err := sqlmock.New()
	if err != nil {
		log.Println("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxMock := sqlx.NewDb(db, "postgres")

	// apply sqlmock obj as global variable
	mock = mockinit

	// Create mock db connection
	err = database.MockDB(sqlxMock, []string{"main"})
	if err != nil {
		log.Println("Fail init mock db connection")
	}

	// Run mini redis
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Create mock redis connection
	cacheConf := make(map[string]string)
	cacheConf["main"] = s.Addr()
	cache.ConnectRedis(cacheConf)
}

func TestNew(t *testing.T) {

	// expect all prepared queries
	prepared = make(map[string]*sqlmock.ExpectedPrepare)
	prepared["get"] = mock.ExpectPrepare("(?i)SELECT id, name, email, phone FROM contacts WHERE id = (.+)")
	prepared["list"] = mock.ExpectPrepare("(?i)SELECT id, name, email, phone FROM contacts ORDER BY id ASC LIMIT (.+) OFFSET (.+)")

	// create pkgcon obj
	pkgCon = New()

	// call it again but make sure it won't prepare queries again
	pkgCon = New()
}

func TestGet(t *testing.T) {

	table := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"phone",
	})

	testCase := []struct {
		ID             int64
		Rows           sqlmock.Rows
		QueryError     bool
		ExpectError    bool
		ExpectedResult Contact
	}{
		{
			1,
			table.AddRow(1, "user1", "user1@email.com", "+628123456789"),
			false,
			false,
			&contact{data: ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"},
		},
		{
			1,
			nil,
			false,
			false,
			&contact{data: ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"},
		},
		{
			2,
			nil,
			true,
			true,
			nil,
		},
	}

	for index, tcase := range testCase {
		if tcase.QueryError {
			prepared["get"].ExpectQuery().WithArgs(tcase.ID).WillReturnError(errors.New("sql error"))
		} else if tcase.Rows != nil {
			prepared["get"].ExpectQuery().WithArgs(tcase.ID).WillReturnRows(tcase.Rows)
		}

		res, err := pkgCon.Get(tcase.ID)
		if (err != nil && !tcase.ExpectError) || (err == nil && tcase.ExpectError) {
			t.Errorf("[TestGet] tcase:%v err got %v | expected err!=nil->%v", index, err, tcase.ExpectError)
		}

		if !reflect.DeepEqual(res, tcase.ExpectedResult) {
			t.Errorf("[TestGet] tcase:%v result got %v | expected %v", index, res, tcase.ExpectedResult)
		}
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestCreate(t *testing.T) {
	table := sqlmock.NewRows([]string{
		"id",
	})

	testCase := []struct {
		Input          ContactData
		Rows           sqlmock.Rows
		QueryError     bool
		ExpectError    bool
		ExpectedResult Contact
	}{
		{
			ContactData{Name: "User1", Email: "user1@email.com", Phone: "+628123456789"},
			table.AddRow(1),
			false,
			false,
			&contact{data: ContactData{ID: 1, Name: "User1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"},
		},
		{
			ContactData{Name: "User2", Email: "user2@email.com", Phone: "+628123456780"},
			nil,
			true,
			true,
			nil,
		},
	}

	for index, tcase := range testCase {
		mock.ExpectBegin()
		query := mock.ExpectQuery("(?i)INSERT INTO contacts name, email, phone VALUES (.+)")
		if tcase.QueryError {
			query.WillReturnError(errors.New("error insert"))
			mock.ExpectRollback()
		} else if tcase.Rows != nil {
			query.WillReturnRows(tcase.Rows)
			mock.ExpectCommit()
		}

		res, err := pkgCon.Create(tcase.Input)
		if !reflect.DeepEqual(res, tcase.ExpectedResult) {
			t.Errorf("[TestCreate] tcase:%v res got %v | expected %v", index, res, tcase.ExpectedResult)
		}

		if (err != nil && !tcase.ExpectError) || (err == nil && tcase.ExpectError) {
			t.Errorf("[TestCreate] tcase:%v err got %v | expected err!=nil->%v", index, err, tcase.ExpectError)
		}

	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestList(t *testing.T) {
	table := sqlmock.NewRows([]string{
		"id",
		"name",
		"email",
		"phone",
	})

	testCase := []struct {
		Take           int64
		Page           int64
		Rows           sqlmock.Rows
		QueryError     bool
		ExpectedResult []ContactData
		ExpectError    bool
	}{
		{
			5,
			1,
			table.AddRow(1, "user1", "user1@email.com", "+628123456789").
				AddRow(2, "user2", "user2@email.com", "+628123456788").
				AddRow(3, "user3", "user3@email.com", "+628123456787").
				AddRow(4, "user4", "user4@email.com", "+628123456786").
				AddRow(5, "user5", "user5@email.com", "+628123456785"),
			false,
			[]ContactData{
				ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"},
				ContactData{ID: 2, Name: "user2", Email: "user2@email.com", Phone: "+628123456788"},
				ContactData{ID: 3, Name: "user3", Email: "user3@email.com", Phone: "+628123456787"},
				ContactData{ID: 4, Name: "user4", Email: "user4@email.com", Phone: "+628123456786"},
				ContactData{ID: 5, Name: "user5", Email: "user5@email.com", Phone: "+628123456785"},
			},
			false,
		},
		{
			10,
			1,
			nil,
			true,
			[]ContactData{},
			true,
		},
		{
			10,
			0,
			nil,
			false,
			[]ContactData{},
			true,
		},
		{
			0,
			3,
			nil,
			false,
			[]ContactData{},
			true,
		},
	}

	for index, tcase := range testCase {
		if tcase.QueryError {
			prepared["list"].ExpectQuery().WillReturnError(errors.New("error get list"))
		} else if tcase.Rows != nil {
			prepared["list"].ExpectQuery().WillReturnRows(tcase.Rows)
		}

		res, err := pkgCon.List(tcase.Take, tcase.Page)

		if !reflect.DeepEqual(res, tcase.ExpectedResult) {
			t.Errorf("[TestList] tcase:%v res got %v | expected %v", index, res, tcase.ExpectedResult)
		}

		if (err != nil && !tcase.ExpectError) || (err == nil && tcase.ExpectError) {
			t.Errorf("[TestList] tcase:%v err got %v | expected err!=nil->%v", index, err, tcase.ExpectError)
		}
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestData(t *testing.T) {
	cObj := &contact{data: ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"}
	data := cObj.Data()
	if !reflect.DeepEqual(data, cObj.data) {
		t.Errorf("[TestData] got %v | expect %v", data, cObj.data)
	}
}

func TestUpdate(t *testing.T) {
	cObj := &contact{data: ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"}

	testCase := []struct {
		Cobj        Contact
		Input       ContactData
		ExpectError bool
		ExpectQuery bool
		QueryArgs   ContactData
	}{
		{
			cObj,
			ContactData{Name: "NewUser1"},
			false,
			true,
			ContactData{Name: "NewUser1", ID: 1, Email: "user1@email.com", Phone: "+628123456789"},
		},
		{
			cObj,
			ContactData{},
			true,
			false,
			ContactData{},
		},
	}

	for index, tcase := range testCase {

		if tcase.ExpectQuery {
			mock.ExpectBegin()
			mock.ExpectExec("(?i)UPDATE contacts SET (.+) WHERE id (.+)").WithArgs(tcase.QueryArgs.Name, tcase.QueryArgs.Email, tcase.QueryArgs.Phone, tcase.QueryArgs.ID).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
		}

		err := cObj.Update(tcase.Input)

		// validate expect error
		if (err != nil && !tcase.ExpectError) || (err == nil && tcase.ExpectError) {
			t.Errorf("[TestUpdate] tcase:%v err got %v | expected err!=nil->%v", index, err, tcase.ExpectError)
		}
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestDelete(t *testing.T) {
	cObj := &contact{data: ContactData{ID: 1, Name: "user1", Email: "user1@email.com", Phone: "+628123456789"}, cacheKey: "contact:1"}

	mock.ExpectBegin()
	mock.ExpectExec("(?i)DELETE FROM contacts WHERE id =(.+)").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := cObj.Delete()
	if err != nil {
		t.Errorf("[TestDelete] fail to delete")
	}
}

func TestValidate(t *testing.T) {
	testCase := []struct {
		Data           ContactData
		ExpectedResult bool
	}{
		{
			ContactData{Name: "Alvin Antonius", Phone: "+628123456789", Email: "alvin.antonius@gmail.com"},
			true,
		},
		{
			ContactData{Name: "Alvin Antonius", Phone: "678+9", Email: "alvin.antonius@gma"},
			false,
		},
		{
			ContactData{Name: "Ad", Phone: "+628123456789", Email: "alvin.antonius@gmail.com"},
			false,
		},
		{
			ContactData{Name: "sdfkh  23", Phone: "+628123456789", Email: "alvin.antonius@gmail.com"},
			true,
		},
		{
			ContactData{Name: "", Phone: "", Email: ""},
			false,
		},
	}

	for index, tcase := range testCase {
		res := validateContact(tcase.Data)
		if res != tcase.ExpectedResult {
			t.Errorf("[TestValidate] tcase:%v res got %v | expected %v", index, res, tcase.ExpectedResult)
		}
	}
}
