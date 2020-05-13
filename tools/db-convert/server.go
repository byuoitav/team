package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/byuoitav/common/db/couch"
)

type device struct {
	ID      string                 `json:"_id"`
	TypeID  string                 `json:"typeID"`
	Address string                 `json:"address"`
	Proxy   map[string]string      `json:"proxy"`
	Ports   []port                 `json:"ports"`
	Tags    map[string]interface{} `json:"tags"`
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

// func updateRooms(db *CouchDB) error {
// 	rooms, err := db.GetAllRooms()
// 	if err != nil {
// 		return fmt.Errorf("error getting all rooms: %s", err)
// 	}

// 	for i := range rooms {
// 		room := rooms[i]
// 		newRoom := structs.Room{
// 			ID:          room.ID,
// 			Designation: room.Designation,
// 			Tags:        room.Tags,
// 		}

// 		_, err := db.UpdateRoom(room.ID, newRoom)
// 		if err != nil {
// 			fmt.Printf("there was an error updating %s: %s\n", rooms.ID, err)
// 		}
// 	}

// 	return nil
// }

// func updateDeviceTypes(db *CouchDB) error {
// 	dTypes, err := db.GetAllDeviceTypes()
// 	if err != nil {
// 		return fmt.Errorf("error getting all device types: %s", err)
// 	}

// 	for i := range dTypes {
// 		dType := dTypes[i]
// 		cmds := make(map[string]structs.Command)
// 		for y := range dType.Commands {
// 			cmd := dType.Commands[y]
// 			if cmds[cmd.ID] != nil {
// 				continue
// 			}

// 			cmds[cmd.ID] = structs.Command{
// 				Address: cmd.Microservice.Address + cmd.Endpoint.Path,
// 				Order:   cmd.Priority,
// 			}
// 		}
// 		newDType := structs.DeviceType{
// 			ID:       dType.ID,
// 			Commands: cmds,
// 		}

// 		err := db.DeleteDeviceType(dType.ID)
// 		if err != nil {
// 			fmt.Printf("unable to delete device type %s: %s\n", dType.ID, err)
// 			continue
// 		}
// 		// this is sketchy cause it could delete fine and not recreate...
// 		_, err = db.CreateDeviceType(newDType)
// 		if err != nil {
// 			fmt.Printf("unable to recreate device type %s: %s\n", newDType.ID, err)
// 		}
// 	}

// 	return nil
// }

// func updateDevices(db *CouchDB) error {
// 	devices, err := db.GetAllDevices()
// 	if err != nil {
// 		return fmt.Errorf("error getting all devices: %s", err)
// 	}

// for i := range devices {
// 	device := devices[i]
// 	newDevice := structs.Device{
// 		ID:      device.ID,
// 		Address: device.Address,
// 		TypeID:  device.Type.ID,
// 		Proxy:   device.Proxy,
// 		Ports:   device.Ports,
// 		Tags:    device.Tags,
// 	}

// 	_, err := db.UpdateDevice(device.ID, newDevice)
// 	if err != nil {
// 		fmt.Printf("unable to update device %s: %s\n", device.ID, err)
// 	}
// }

// 	return nil
// }

func updateDevice(db *couch.CouchDB) error {

	room, err := db.GetDevicesByRoom("ITB-1109B")
	if err != nil {
		return fmt.Errorf("error getting room: %s", err)
	}
	newDB := couch.NewDB(os.Getenv("DB_ADDRESS"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))

	for i := range room {
		dev := room[i]
		newDevice := device{
			ID:      dev.ID,
			Address: dev.Address,
			TypeID:  dev.Type.ID,
			Proxy:   dev.Proxy,
		}
		var ports []port
		for y := range dev.Ports {
			p := dev.Ports[y]
			newPort := port{
				Name:     p.ID,
				Endpoint: p.DestinationDevice,
				Type:     p.PortType,
			}

			ports = append(ports, newPort)
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
