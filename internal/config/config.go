package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type (
	// C is stand for config object
	C struct {
		Database map[string]*struct {
			Master string
			Slave  string
		} `json:"database"`

		Redis map[string]string `json:"redis"`

		Port string `json:"port"`
	}

	dbconf struct {
		Master string `json:"master"`
		Slave  string `json:"slave"`
	}
)

var configuration C

// ReadConfig is for read configuration from multiple path
// it'll loop all filepath and stop if file is found
func ReadConfig(filepath ...string) *C {

	for _, path := range filepath {
		// read config file
		fileByte, err := ioutil.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}

		err = json.Unmarshal(fileByte, &configuration)
		if err != nil {
			log.Fatal(err)
		}
		break
	}

	return &configuration
}

// Get config data
func Get() *C {
	return &configuration
}
