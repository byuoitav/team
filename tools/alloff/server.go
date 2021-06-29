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
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/common/v2/events"
)

var version = "0.7"
var hub_url = ""
var layout = "2006-01-02T15:04:05"

// building Publisher in order to send messages to
func Publish(url string, event events.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Unable to Marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Unable to build request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to make request: %w", err)
	}

	defer resp.Body.Close()

	fmt.Printf("Response code: %v Body: %v\n", resp.StatusCode, resp.Body)

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("got a %v response from event url", resp.StatusCode)
	}

	return nil
}

// Log file path building.
func log_file_path() string {
	// Generate filename
	t := time.Now()
	file := t.Format(layout) + "_nightly_restart.log"
	// Generate full path
	dir := "/var/log/restart/"
	path := dir + file

	// Check if Path exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0644)
		if err != nil {
			fmt.Printf("Error Creating File: %w\n", err.Error())
		}
	}

	return path
}
func main() {

	file_path := log_file_path()

	f, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf(err.Error())
	}

	defer f.Close()
	t := time.Now()
	fmt.Fprintf(f, "Start Time: %v\n", t.Format(layout))
	fmt.Fprintf(f, "Version: %v\n", version)

	//regex to extract the room id from the hostname
	re := regexp.MustCompile("([A-Z,0-9]+-[A-Z,0-9]+)-[A-Z,0-9]+")

	host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(f, "failed to get hostname: %s\n", err)
		return
	}
	rm := re.FindStringSubmatch(host)[1]

	hub_url := fmt.Sprintf("http://%v:7100/event", host)

	split := strings.Split(host, "-")

	bld := split[0]
	room_name := split[0] + "-" + split[1]

	//get the list of devices
	db := db.GetDBWithCustomAuth(os.Getenv("COUCH_ADDR"), os.Getenv("COUCH_USER"), os.Getenv("COUCH_PASS"))
	devs, err := db.GetDevicesByRoomAndRole(rm, "VideoOut")
	if err != nil {
		fmt.Fprintf(f, "Error getting Video out devices from database: %v.\n", err.Error())
		return
	}

	// Get the list of Cameras
	cams, err := db.GetDevicesByRoomAndType(rm, "AVER 520 Pro Camera")
	if err != nil {
		fmt.Fprintf(f, "Error Getting Cameras from Database: %v.\n", err.Error())
		return
	}

	room := structs.PublicRoom{}

	for i := range devs {
		if !strings.Contains(devs[i].ID, "-CAM") {
			room.Displays = append(room.Displays, structs.Display{
				PublicDevice: structs.PublicDevice{
					Name:  devs[i].Name,
					Power: "standby",
				},
			})
		} else {
			fmt.Fprintf(f, "%v doesn't support standby, skipping\n", devs[i].Name)
		}
	}

	b, err := json.Marshal(&room)
	if err != nil {
		fmt.Fprintf(f, "Couldn't marshal the room: %v\n", err.Error())
	}

	//make our request
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8000/buildings/%v/rooms/%v", split[0], split[1]), bytes.NewReader(b))
	if err != nil {
		fmt.Fprintf(f, "Couldn't create request: %v\n", err.Error())
	}

	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(f, "Couldn't make request\n", err.Error())
	}

	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(f, "Couldn't read body: %v\n", err.Error())
	}

	rest_resp := fmt.Sprintf("Finished setting room to standby. Response Code: %v. Response Body: %s\n", resp.StatusCode, out)

	event := events.Event{
		GeneratingSystem: host,
		Timestamp:        time.Now(),
		User:             host,
		Key:              "nightly-reboot",
		Value:            "system-standby",
		EventTags: []string{
			"auto-generated",
			"restart",
		},
		Data: rest_resp,
	}

	event.TargetDevice.BuildingID = bld
	event.TargetDevice.RoomID = room_name
	event.TargetDevice.DeviceID = host

	err = Publish(hub_url, event)
	if err != nil {
		fmt.Fprintf(f, "Sending event to event hub failed: %w", err.Error())
	}

	fmt.Fprintf(f, rest_resp)

	// Reboot the Cameras in the room
	if strings.Contains(host, "-CP1") && len(cams) > 0 {
		for _, cam := range cams {
			var camfail string
			var Cam_txt string
			camname := cam.Address
			req, err := http.NewRequest("GET", fmt.Sprintf("https://aver.av.byu.edu/v1/Pro520/%v:52381/reboot", camname), nil)
			if err != nil {
				fmt.Fprintf(f, "Couldn't restart cameras\n")
			}
			req.Header.Add("content-type", "application/json")

			camresp, err := client.Do(req)
			if err != nil {
				camfail := fmt.Sprintf("Couldn't make request: %v\n", err.Error())
				fmt.Fprintf(f, camfail)
			}

			defer camresp.Body.Close()

			if camresp.StatusCode/100 != 2 {
				fmt.Fprintf(f, "Non-200 response received: %v\n", camresp.StatusCode)
				camout, err := ioutil.ReadAll(camresp.Body)
				if err != nil {
					fmt.Fprintf(f, "Cannot read the response body", err.Error())
				}
				fmt.Fprintf(f, "Response Body: %s\n", camout)

			}

			if camresp != nil {
				Cam_txt = fmt.Sprintf("Finished rebooting %v. Response Code: %v\n", camname, camresp.StatusCode)
			} else if camfail != "" {
				Cam_txt = fmt.Sprintf("Camera Reboot Failed for %v. Error: %v\n", camname, camfail)
			}

			event := events.Event{
				GeneratingSystem: host,
				Timestamp:        time.Now(),
				User:             host,
				Key:              "nightly-reboot",
				Value:            "camera-reboot",
				EventTags: []string{
					"auto-generated",
					"restart",
				},
				Data: Cam_txt,
			}

			event.TargetDevice.BuildingID = bld
			event.TargetDevice.RoomID = room_name
			event.TargetDevice.DeviceID = camname

			err = Publish(hub_url, event)
			if err != nil {
				fmt.Fprintf(f, "Could not publish event to the event hub: %w\n", err.Error())
			}
			fmt.Fprintf(f, Cam_txt)
		}
	} else {
		fmt.Fprintln(f, "No cameras in this room to reboot.")
	}

	// Refresh UI
	fmt.Fprintln(f, "Refreshing Pi UI")

	client = http.Client{
		Timeout: 5 * time.Second,
	}
	req, err = http.NewRequest("PUT", fmt.Sprintf("http://localhost:8888/refresh"), nil)
	if err != nil {
		fmt.Fprintf(f, "Couldn't create request: %v\n", err.Error())
	}

	req.Header.Add("content-type", "application/json")

	resp1, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(f, "Couldn't make request\n", err.Error())
	}
	if resp1.StatusCode/100 != 2 {
		fmt.Fprintf(f, "Non-200 received\n", err.Error())
	}

	defer resp1.Body.Close()
	fmt.Fprintf(f, "done refreshing..\n")

}
