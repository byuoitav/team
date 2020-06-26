package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/common/db/couch"
)

type device struct {
	ID      string                 `json:"_id"`
	TypeID  string                 `json:"typeID"`
	Address string                 `json:"address,omitempty"`
	Proxy   map[string]string      `json:"proxy,omitempty"`
	Ports   []port                 `json:"ports,omitempty"`
	Tags    map[string]interface{} `json:"tags,omitempty"`
}

type TypeID struct {
	ID string `json:"_id"`
}

type deviceType struct {
	ID       string             `json:"_id"`
	Commands map[string]command `json:"commands"`
}

type command struct {
	URLs  map[string]string `json:"addresses"`
	Order *int              `json:"order,omitempty"`
}

type port struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Incoming bool   `json:"incoming"`
	Outgoing bool   `json:"outgoing"`
	Type     string `json:"type"`
}

type CouchUpsertResponse struct {
	OK  bool   `json:"ok"`
	ID  string `json:"id"`
	Rev string `json:"rev"`
}

type Room struct {
	ID          string                 `json:"_id"`
	Designation string                 `json:"designation"`
	Tags        map[string]interface{} `json:"tags,omitempty"`
}

func main() {
	db := couch.NewDB(os.Getenv("DB_ADDRESS"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	err := updateDevice(db)
	if err != nil {
		fmt.Println("poop: ", err)
	}
	// err := updateRooms(db)
	// if err != nil {
	// 	fmt.Printf("unable to update rooms: %s\n", err)
	// }

	// err = updateDeviceTypes(db)
	// if err != nil {
	// 	fmt.Printf("unable to update device types: %s\n", err)
	// }

	// err = updateDevices(db)
	// if err != nil {
	// 	fmt.Printf("unable to update devices: %s\n", err)
	// }
}

// func updateRooms(db *couch.CouchDB) error {
// 	rooms, err := db.GetAlls()
// 	if err != nil {
// 		return fmt.Errorf("error getting all rooms: %s", err)
// 	}

// 	for i := range rooms {
// 		room := rooms[i]
// 		newRoom := Room{
// 			ID:          room.ID,
// 			Designation: room.Designation,
// 		}

// 		if room.Description != "" {
// 			newRoom.Tags = map[string]interface{}{
// 				"notes": room.Description,
// 			}
// 		}

// 		if err := db.DeleteRoom(room.ID); err != nil {
// 			fmt.Printf("error deleting room %s: %s\n", room.ID, err)
// 			continue
// 		}

// 		b, err := json.Marshal(newRoom)
// 		if err != nil {
// 			fmt.Printf("error marshalling room %s into json: %s\n", newRoom.ID, err)
// 			continue
// 		}

// 		var resp CouchUpsertResponse
// 		err = db.MakeRequest("POST", "rooms", "application/json", b, &resp)
// 		if err != nil {
// 			fmt.Printf("error recreating room %s: %s\n", newRoom.ID, err)
// 		}
// 	}

// 	return nil
// }

// func updateDeviceTypes(db *couch.CouchDB) error {
// 	dTypes, err := db.GetAllDeviceTypes()
// 	if err != nil {
// 		return fmt.Errorf("error getting all device types: %s", err)
// 	}

// 	for i := range dTypes {
// 		dType := dTypes[i]
// 		cmds := make(map[string]command)
// 		for y := range dType.Commands {
// 			cmd := dType.Commands[y]
// 			if cmds[cmd.ID] != nil {
// 				continue
// 			}

// 			//parse the path and replace :variable with {{variable}}
// 			addr := cmd.Microservice.Address
// 			splitPath := strings.Split(cmd.Endpoint.Path, "/")
// 			for x := range splitPath {
// 				if splitPath[x][0] == ':' {
// 					splitPath[x] = strings.TrimPrefix(splitPath[x], ":")
// 					splitPath[x] = "{{" + splitPath[x] + "}}"
// 				}
// 				addr = addr + "/" + splitPath[x]
// 			}

// 			cmds[cmd.ID] = command{
// 				URLs: map[string]string{
// 					"fallback": addr,
// 				},
// 				Order: cmd.Priority,
// 			}
// 		}
// 		newDType := deviceType{
// 			ID:       dType.ID,
// 			Commands: cmds,
// 		}

// 		// this is sketchy cause it could delete fine and not recreate...

// 		err := db.DeleteDeviceType(dType.ID)
// 		if err != nil {
// 			fmt.Printf("unable to delete device type %s: %s\n", dType.ID, err)
// 			continue
// 		}

// 		b, err := json.Marshal(newDType)
// 		if err != nil {
// 			fmt.Printf("error marshalling device type %s into json: %s\n", newDType.ID, err)
// 			continue
// 		}

// 		var resp CouchUpsertResponse
// 		err = db.MakeRequest("POST", "device-types", "application/json", b, &resp)
// 		if err != nil {
// 			fmt.Printf("error recreating device type %s: %s\n", newDType.ID, err)
// 		}
// 	}

// 	return nil
// }

// func updateDevices(db *couch.CouchDB) error {
// 	devices, err := db.GetAllDevices()
// 	if err != nil {
// 		return fmt.Errorf("error getting all devices: %s", err)
// 	}

// 	portMap := make(map[string]string)
// 	for i := range devices {
// 		dev := devices[i]
// 		for x := range dev.Ports {
// 			p := dev.Ports[x]
// 			for y := range p.Tags {
// 				t := p.Tags[y]
// 				if t == "port-in" {
// 					portMap[p.SourceDevice] = p.DestinationDevice
// 				}
// 			}
// 		}
// 	}

// 	for i := range devices {
// 		dev := devices[i]
// 		newDevice := device{
// 			ID:     dev.ID,
// 			TypeID: dev.Type.ID,
// 			Proxy:  dev.Proxy,
// 			Tags:   dev.Tags,
// 		}
// 		var ports []port
// 		for y := range dev.Ports {
// 			p := dev.Ports[y]
// 			newPort := port{
// 				Name: p.ID,
// 				Type: "video",
// 			}

// 			if strings.Contains(newPort.Name, "IN") {
// 				newPort.Name = strings.TrimPrefix(newPort.Name, "IN")
// 				newPort.Outgoing = false
// 				newPort.Incoming = true
// 			}
// 			if strings.Contains(newPort.Name, "OUT") {
// 				newPort.Name = strings.TrimPrefix(newPort.Name, "OUT")
// 				newPort.Incoming = false
// 				newPort.Outgoing = true
// 			}

// 			if newDevice.Address != "0.0.0.0" {
// 				newDevice.Address = dev.Address
// 			}

// 			audio := false
// 			video := false

// 			for x := range p.Tags {
// 				if p.Tags[x] == "audio" {
// 					audio = true
// 				}
// 				if p.Tags[x] == "video" {
// 					video = true
// 				}
// 			}
// 			if audio && !video {
// 				newPort.Type = "audio"
// 			}
// 			tid := newDevice.TypeID

// 			if tid == "SonyXBR" || tid == "Sony PHZ10" || tid == "Sony Projector No Audio" || tid == "Sony XBR No Audio" ||
// 				tid == "Sony XBR Nop Audio (Mirror)" {
// 				newPort.Outgoing = false
// 				newPort.Incoming = true
// 			}
// 			if newDevice.TypeID == "QSC-Core-110F" {
// 				newPort.Outgoing = false
// 				newPort.Incoming = true
// 				newPort.Type = "audio"
// 			}

// 			if newDevice.ID == p.DestinationDevice {
// 				newPort.Outgoing = false
// 				newPort.Incoming = true
// 				newPort.Endpoint = p.SourceDevice
// 			} else if newDevice.ID == p.SourceDevice {
// 				newPort.Outgoing = true
// 				newPort.Incoming = false
// 				newPort.Endpoint = p.DestinationDevice
// 			}

// 			ports = append(ports, newPort)
// 		}

// 		for y := range dev.Roles {
// 			if dev.Roles[y].ID == "VideoIn" {
// 				newPort := port{
// 					Endpoint: portMap[dev.ID],
// 					Outgoing: true,
// 					Type:     "video",
// 				}
// 				ports = append(ports, newPort)
// 			}
// 		}
// 		newDevice.Ports = ports

// 		if dev.Description != "" {
// 			if newDevice.Tags == nil {
// 				newDevice.Tags = map[string]interface{}{
// 					"notes": dev.Description,
// 				}
// 			} else {
// 				newDevice.Tags["description"] = dev.Description
// 			}
// 		}

// 		if err = db.DeleteDevice(dev.ID); err != nil {
// 			fmt.Printf("error deleting device %s: %s", dev.ID, err)
// 			continue
// 		}

// 		// re-add the device
// 		b, err := json.Marshal(newDevice)
// 		if err != nil {
// 			fmt.Printf("There was an error marshaling the device: %s\n", err)
// 			continue
// 		}

// 		var resp CouchUpsertResponse
// 		err = db.MakeRequest("POST", fmt.Sprintf("%v", "devices"), "application/json", b, &resp)
// 		if err != nil {
// 			fmt.Printf("error making request: %s\n", err)
// 			continue
// 		}
// 	}

// 	return nil
// }

//this actually does a room lol
func updateDevice(db *couch.CouchDB) error {

	room, err := db.GetDevicesByRoom("ITB-1108A")
	if err != nil {
		return fmt.Errorf("error getting room: %s", err)
	}
	newDB := couch.NewDB("https://couchdb-dev.avs.byu.edu", "dev", os.Getenv("DB_PASSWORD"))

	portMap := make(map[string]string)
	for i := range room {
		dev := room[i]
		for x := range dev.Ports {
			p := dev.Ports[x]
			if strings.Contains(p.SourceDevice, "MIC") {
				portMap[p.SourceDevice] = p.DestinationDevice
			}
			for y := range p.Tags {
				t := p.Tags[y]
				if t == "port-in" {
					portMap[p.SourceDevice] = p.DestinationDevice
				}
			}
		}
	}

	for i := range room {
		dev := room[i]
		newDevice := device{
			ID:     dev.ID,
			TypeID: dev.Type.ID,
			Proxy:  dev.Proxy,
		}
		var ports []port
		for y := range dev.Ports {
			p := dev.Ports[y]
			newPort := port{
				Name: p.ID,
				Type: "video",
			}

			if strings.Contains(newPort.Name, "IN") {
				newPort.Name = strings.TrimPrefix(newPort.Name, "IN")
				newPort.Outgoing = false
				newPort.Incoming = true
			}
			if strings.Contains(newPort.Name, "OUT") {
				newPort.Name = strings.TrimPrefix(newPort.Name, "OUT")
				newPort.Incoming = false
				newPort.Outgoing = true
			}

			if newDevice.Address != "0.0.0.0" {
				newDevice.Address = dev.Address
			}

			// default type is video, if audio then audio, if it has both audio and video, then video
			audio := false
			video := false

			for x := range p.Tags {
				if p.Tags[x] == "audio" {
					audio = true
				}
				if p.Tags[x] == "video" {
					video = true
				}
			}
			if audio && !video {
				newPort.Type = "audio"
			}
			tid := newDevice.TypeID

			if tid == "SonyXBR" || tid == "Sony PHZ10" || tid == "Sony Projector No Audio" || tid == "Sony XBR No Audio" ||
				tid == "Sony XBR Nop Audio (Mirror)" {
				newPort.Outgoing = false
				newPort.Incoming = true
			}
			if newDevice.TypeID == "QSC-Core-110F" {
				newPort.Outgoing = false
				newPort.Incoming = true
				newPort.Type = "audio"
			}

			if newDevice.ID == p.DestinationDevice {
				newPort.Outgoing = false
				newPort.Incoming = true
				newPort.Endpoint = p.SourceDevice
			} else if newDevice.ID == p.SourceDevice {
				newPort.Outgoing = true
				newPort.Incoming = false
				newPort.Endpoint = p.DestinationDevice
			}

			ports = append(ports, newPort)
		}

		if strings.Contains(dev.ID, "MIC") {
			newPort := port{
				Endpoint: portMap[dev.ID],
				Outgoing: true,
				Type:     "audio",
			}
			ports = append(ports, newPort)
		}

		for y := range dev.Roles {
			if dev.Roles[y].ID == "VideoIn" {
				newPort := port{
					Endpoint: portMap[dev.ID],
					Outgoing: true,
					Type:     "video",
				}
				ports = append(ports, newPort)
				break
			}
		}
		newDevice.Ports = ports

		b, err := json.Marshal(newDevice)
		if err != nil {
			fmt.Printf("There was an error marshaling the device: %s\n", err)
			continue
		}

		var resp CouchUpsertResponse
		err = newDB.MakeRequest("POST", fmt.Sprintf("%v", "devices"), "application/json", b, &resp)
		if err != nil {
			fmt.Printf("error making request: %s\n", err)
			continue
		}
	}

	return nil
}
