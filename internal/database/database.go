package database

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type (
	databaseReplication struct {
		Master *sqlx.DB
		Slave  *sqlx.DB
	}
)

var databases map[string]databaseReplication

// ConnectDB for create db connection
func ConnectDB(config interface{}) {
	cfg := config.(map[string]*struct {
		Master string
		Slave  string
	})

	// create map for store all db conn
	databases = make(map[string]databaseReplication)

	for dbname, conn := range cfg {
		// connect to master DB
		dbmaster, err := sqlx.Connect("postgres", conn.Master)
		if err != nil {
			log.Fatalln(err)
		}
		// these are just my number for limit db open conn
		dbmaster.SetMaxIdleConns(3)
		dbmaster.SetMaxOpenConns(10)

		// connect to slave DB
		dbslave, err := sqlx.Connect("postgres", conn.Slave)
		if err != nil {
			log.Fatalln(err)
		}
		// these are just my number for limit db open conn
		dbslave.SetMaxIdleConns(3)
		dbslave.SetMaxOpenConns(10)

		// assign db conn to struct
		// assign struct to map
		databases[dbname] = databaseReplication{
			Master: dbmaster,
			Slave:  dbslave,
		}
	}
}

// Conn is for get database connection
func Conn(dbname, replication string) (*sqlx.DB, error) {
	if dbconn, ok := databases[dbname]; ok {
		if replication == "master" {
			return dbconn.Master, nil
		}
		return dbconn.Slave, nil
	}

	return nil, fmt.Errorf("database %s not found", dbname)
}

// MockDB is for unit testing that require mocking DB
func MockDB(mockdb *sqlx.DB, replications []string) error {

	if len(replications) == 0 {
		return fmt.Errorf("no database replication listed")
	}

	databases = make(map[string]databaseReplication)
	for _, repl := range replications {
		databases[repl] = databaseReplication{
			Master: mockdb,
			Slave:  mockdb,
		}
	}

	return nil
}
