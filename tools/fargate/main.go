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
	newService := flag.Bool("n", true, "Create a new cluster, service, load balancer, etc.")

	DBName := flag.String("dbname", "deployment-information", "The databse name for service configuration info")
	Output := flag.String("o", "", "If defined will output the template generated to this file. Will not attempt to create/update the stack.")

	flag.Parse()

	log.L.Infof("%v", *newService)
	log.L.Infof("%v", *Output)

	//we need to go get the config file
	configdef, err := ReadConfigFile(*config)
	if err != nil {
		return
	}
	var b []byte
	var fileName string

	//we need to go get the datbase information for this service
	configwrap, err := GetInfoFromDB(*DBName, configdef.Name)
	if err != nil {
		log.L.Fatalf(err.Error())
	}

	if *newService {
		//we need to create the cluster, etc. etc.a
		config, name, err := BuildNewService(configwrap, configdef, *DBName, *branch)

		if err != nil {
			log.L.Fatalf(err.Error())
		}

		b, err = json.MarshalIndent(config, "", " ")
		if err != nil {
			log.L.Fatalf(err.Error())
		}
		fileName = name
	} else {

		taskConfig, name, err := buildTaskDefinitionConfig(configwrap, configdef, *DBName, *branch)
		if err != nil {
			log.L.Fatalf(err.Error())
		}

		b, err = json.MarshalIndent(taskConfig, "", " ")
		if err != nil {
			log.L.Fatalf(err.Error())
		}
		fileName = name

	}
	if len(*Output) > 0 {
		log.L.Infof("Printing to file %v", *Output)
		//we output to the file specified
		err := ioutil.WriteFile(*Output, b, 0644)
		if err != nil {
			log.L.Errorf("Couldn't write the file: %v", err.Error())
			return
		}
	} else {
		err := StartDeployment(fileName, b)
		if err != nil {
			log.L.Errorf("Couldn't start cloudformation deployment: %v", err.Error())
			return
		}
	}
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
func GetInfoFromDB(dbname, service string) (ConfigInfoWrapper, error) {
	var toReturn ConfigInfoWrapper
	v, ok := db.GetDB().(*couch.CouchDB)
	if !ok {
		return toReturn, errors.New("unkown database type")
	}

	err := v.MakeRequest("GET", fmt.Sprintf("%v/%v", dbname, service), "application/json", []byte{}, &toReturn)
	if err != nil {
		log.L.Errorf("%v", err.Error())
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
		log.L.Errorf("%v", err.Error())
		return toReturn, fmt.Errorf("Couldn't get task info for %v from database: %v", taskname, err.Error())
	}
	return toReturn, nil
}
