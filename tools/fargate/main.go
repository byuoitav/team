package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/db/couch"
	"github.com/byuoitav/common/log"
)

func main() {
	branch := flag.String("b", "development", "the branch to deploy to")
	config := flag.String("c", "./config.json", "the location of the config file")

	DBName := flag.String("dbname", "deployment-information", "The user to log into the database with")

	//we need to go get the config file
	configdef, err := ReadConfigFile(*config)
	if err != nil {
		return
	}

	//we need to go get the datbase information for this service
	configwrap, err := GetInfoFromDB(*DBName, configdef.Name)
	if err != nil {
		log.L.Fatalf(err.Error())
	}

	taskConfig, err := buildAWSConfig(configwrap, configdef, *DBName, *branch)
	if err != nil {
		log.L.Fatalf(err.Error())
	}

	b, err := json.MarshalIndent(taskConfig, "", " ")
	if err != nil {
		log.L.Fatalf(err.Error())
	}
	fmt.Printf("%s", b)
}

//ReadConfigFile .
func ReadConfigFile(a string) (ConfigDefinition, error) {
	var toReturn ConfigDefinition

	b, err := ioutil.ReadFile(a)
	if err != nil {
		log.L.Errorf("Couldn't read file %v, %v", a, err.Error())
		return toReturn, err
	}

	err = json.Unmarshal(b, &toReturn)
	if err != nil {
		log.L.Errorf("Invalid config: %v", err.Error())
		return toReturn, err
	}

	return toReturn, nil
}

//GetInfoFromDB .
func GetInfoFromDB(name, service string) (ConfigInfoWrapper, error) {
	var toReturn ConfigInfoWrapper
	v, ok := db.GetDB().(*couch.CouchDB)
	if !ok {
		return toReturn, errors.New("unkown database type")
	}

	err := v.MakeRequest("GET", fmt.Sprintf("%v/%v", name, service), "application/json", []byte{}, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("Couldn't get config info from database: %v", err.Error())
	}

	return toReturn, nil
}

//GetTaskInfoFromDB .
func GetTaskInfoFromDB(taskname string) (AWSTaskWrapper, error) {
	var toReturn AWSTaskWrapper

	v, ok := db.GetDB().(*couch.CouchDB)
	if !ok {
		return toReturn, errors.New("unkown database type")
	}

	err := v.MakeRequest("GET", fmt.Sprintf("aws-deployment-info/%v", taskname), "application/json", []byte{}, &toReturn)
	if err != nil {
		return toReturn, fmt.Errorf("Couldn't get config info from database: %v", err.Error())
	}
	return toReturn, nil
}
