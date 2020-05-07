package main

import (
	"fmt"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/structs"
)

func main() {
	db := db.GetDBWithCustomAuth("http://couch-dev.avs.byu.edu", "dev", "password")
	err := updateRooms(db)
	if err != nil {
		fmt.Printf("unable to update rooms: %s\n", err)
	}

	err = updateDeviceTypes(db)
	if err != nil {
		fmt.Printf("unable to update device types: %s\n", err)
	}

	err = updateDevices(db)
	if err != nil {
		fmt.Printf("unable to update devices: %s\n", err)
	}
}

func updateRooms(db *CouchDB) error {
	rooms, err := db.GetAllRooms()
	if err != nil {
		return fmt.Errorf("error getting all rooms: %s", err)
	}

	for i := range rooms {
		room := rooms[i]
		newRoom := structs.Room{
			ID:          room.ID,
			Designation: room.Designation,
			Tags:        room.Tags,
		}

		_, err := db.UpdateRoom(room.ID, newRoom)
		if err != nil {
			fmt.Printf("there was an error updating %s: %s\n", rooms.ID, err)
		}
	}

	return nil
}

func updateDeviceTypes(db *CouchDB) error {
	dTypes, err := db.GetAllDeviceTypes()
	if err != nil {
		return fmt.Errorf("error getting all device types: %s", err)
	}

	for i := range dTypes {
		dType := dTypes[i]
		cmds := make(map[string]structs.Command)
		for y := range dType.Commands {
			cmd := dType.Commands[y]
			if cmds[cmd.ID] != nil {
				continue
			}

			cmds[cmd.ID] = structs.Command{
				Address: cmd.Microservice.Address + cmd.Endpoint.Path,
				Order:   cmd.Priority,
			}
		}
		newDType := structs.DeviceType{
			ID:       dType.ID,
			Commands: cmds,
		}

		err := db.DeleteDeviceType(dType.ID)
		if err != nil {
			fmt.Printf("unable to delete device type %s: %s\n", dType.ID, err)
			continue
		}
		// this is sketchy cause it could delete fine and not recreate...
		_, err = db.CreateDeviceType(newDType)
		if err != nil {
			fmt.Printf("unable to recreate device type %s: %s\n", newDType.ID, err)
		}
	}

	return nil
}

func updateDevices(db *CouchDB) error {
	devices, err := db.GetAllDevices()
	if err != nil {
		return fmt.Errorf("error getting all devices: %s", err)
	}

	for i := range devices {
		device := devices[i]
		newDevice := structs.Device{
			ID:      device.ID,
			Address: device.Address,
			TypeID:  device.Type.ID,
			Proxy:   device.Proxy,
			Ports:   device.Ports,
			Tags:    device.Tags,
		}

		_, err := db.UpdateDevice(device.ID, newDevice)
		if err != nil {
			fmt.Printf("unable to update device %s: %s\n", device.ID, err)
		}
	}

	return nil
}
