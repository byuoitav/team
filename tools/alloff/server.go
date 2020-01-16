package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
)

var version = "0.5"

func main() {
	log.L.Infof("Version: %v", version)

	//regex to extract the room id from the hostname
	re := regexp.MustCompile("([A-Z,0-9]+-[A-Z,0-9]+)-[A-Z,0-9]+")

	host, err := os.Hostname()
	if err != nil {
		log.L.Errorf("failed to get hostname: %s", err)
		return
	}

	rm := re.FindStringSubmatch(host)[1]

	split := strings.Split(host, "-")

	//get the list of devices
	db := db.GetDBWithCustomAuth(os.Getenv("COUCH_ADDR"), os.Getenv("COUCH_USER"), os.Getenv("COUCH_PASS"))
	devs, err := db.GetDevicesByRoomAndRole(rm, "VideoOut")
	if err != nil {
		log.L.Errorf("Couldn't get Video out devices in the room to turn off: %v.", err.Error())
		return
	}

	room := structs.PublicRoom{}

	for i := range devs {
		room.Displays = append(room.Displays, structs.Display{
			PublicDevice: structs.PublicDevice{
				Name:  devs[i].Name,
				Power: "standby",
			},
		})
	}

	b, err := json.Marshal(&room)
	if err != nil {
		log.L.Fatalf("Couldn't marshal the room: %v", err.Error())
	}

	//make our request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8000/buildings/%v/rooms/%v", split[0], split[1]), bytes.NewReader(b))
	if err != nil {
		log.L.Fatalf("Couldn't create request: %v", err.Error())
	}

	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.L.Fatalf("Couldn't make request", err.Error())
	}

	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Fatalf("Couldn't read body: %v", err.Error())
	}

	log.L.Infof("Finished running. Response Code: %v. Response Body: %s", resp.StatusCode, out)
	log.L.Infof("Refreshing")

	client = http.Client{
		Timeout: 5 * time.Second,
	}
	req, err = http.NewRequest("PUT", fmt.Sprintf("http://localhost:8888/refresh"), nil)
	if err != nil {
		log.L.Fatalf("Couldn't create request: %v", err.Error())
	}

	req.Header.Add("content-type", "application/json")

	resp1, err := client.Do(req)
	if err != nil {
		log.L.Fatalf("Couldn't make request", err.Error())
	}
	if resp1.StatusCode/100 != 2 {
		log.L.Fatalf("Non-200 received", err.Error())
	}

	defer resp1.Body.Close()
	log.L.Infof("done refreshing..")

}
