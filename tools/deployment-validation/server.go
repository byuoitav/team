package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/status"
	"github.com/byuoitav/common/structs"
)

// Report .
type Report struct {
	ID      string
	Good    bool
	Message string
}

func main() {
	log.SetLevel("info")

	devices, gerr := db.GetDB().GetDevicesByRoleAndType("ControlProcessor", "Pi3")
	if gerr != nil {
		log.L.Fatalf("failed: %s", gerr.Error())
	}

	versions := map[string]string{
		"device-monitoring": "0.0.132",
		"central-event-hub": "0.0.2",
		"ui":                "1.0.13",
		"av-api":            "0.1.93",
	}

	total := 0
	done := 0
	reports := make(chan Report, 1000)
	wg := sync.WaitGroup{}

	for _, device := range devices {
		if strings.Contains(device.ID, "MONK") ||
			strings.Contains(device.ID, "STAGE") ||
			strings.Contains(device.ID, "HEALTH") ||
			strings.Contains(device.ID, "DNB") ||
			strings.Contains(device.ID, "KEN") ||
			strings.Contains(device.ID, "SWKT") ||
			strings.Contains(device.ID, "BOB") ||
			strings.Contains(device.ID, "ITB") ||
			strings.Contains(device.ID, "JOHN") ||
			strings.Contains(device.ID, "CCUMC") ||
			strings.Contains(device.ID, "TEST") {
			continue
		}

		total++
	}

	for _, device := range devices {
		if strings.Contains(device.ID, "MONK") ||
			strings.Contains(device.ID, "STAGE") ||
			strings.Contains(device.ID, "HEALTH") ||
			strings.Contains(device.ID, "DNB") ||
			strings.Contains(device.ID, "KEN") ||
			strings.Contains(device.ID, "SWKT") ||
			strings.Contains(device.ID, "BOB") ||
			strings.Contains(device.ID, "ITB") ||
			strings.Contains(device.ID, "JOHN") ||
			strings.Contains(device.ID, "CCUMC") ||
			strings.Contains(device.ID, "TEST") {
			continue
		}

		wg.Add(1)

		go func(device structs.Device) {
			url := fmt.Sprintf("http://%s:10000/device/status", device.Address)
			report := Report{
				ID:   device.ID,
				Good: false,
			}

			defer func() {
				reports <- report
				done++

				log.L.Infof("finished %d of %d", done, total)
				wg.Done()
			}()

			if strings.Contains(url, "0.0.0.0") {
				report.Message = fmt.Sprintf("invalid address: %s", device.Address)
				return
			}

			client := &http.Client{
				Timeout: 10 * time.Second,
			}

			resp, err := client.Get(url)
			if err != nil {
				report.Message = fmt.Sprintf("unable to make GET request against %s: %s", url, err)
				return
			}
			defer resp.Body.Close()

			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				report.Message = fmt.Sprintf("unable to read response body: %s", err)
				return
			}

			stati := []status.Status{}
			err = json.Unmarshal(bytes, &stati)
			if err != nil {
				report.Message = fmt.Sprintf("unable to unmarshal response: %s", err)
				return
			}

			// check if the response is good
			for name, version := range versions {
				idx := -1

				for i := range stati {
					if stati[i].Name == name {
						idx = i
						break
					}
				}

				if idx == -1 {
					report.Message = fmt.Sprintf("missing microservice in response: %s", name)
					return
				}

				if stati[idx].Version != version {
					report.Message = fmt.Sprintf("incorrect version for %s: expected %s, got %s", name, version, stati[idx].Version)
					return
				}
			}

			report.Good = true

			// hit endpoint(s)
			url = fmt.Sprintf("http://%s:10000/room/hardwareinfo", device.Address)
			_, err = client.Get(url)
			if err != nil {
				report.Message = fmt.Sprintf("failed to get room hardware info: %s", err)
				return
			}

			url = fmt.Sprintf("http://%s:10000/room/state", device.Address)
			_, err = client.Get(url)
			if err != nil {
				report.Message = fmt.Sprintf("failed to get room state: %s", err)
				return
			}

			url = fmt.Sprintf("http://%s:10000/device/hardwareinfo", device.Address)
			_, err = client.Get(url)
			if err != nil {
				report.Message = fmt.Sprintf("failed to get device hardware info: %s", err)
				return
			}

			report.Message = fmt.Sprintf("successfully sent hardware info up")
		}(device)
	}

	wg.Wait()
	close(reports)

	compiled := []Report{}
	for report := range reports {
		compiled = append(compiled, report)
	}

	f, err := os.Create("reports.txt")
	if err != nil {
		log.L.Fatalf("failed to open file for writing", err.Error())
	}

	defer func() {
		f.Sync()
		f.Close()
	}()

	_, err = f.Write([]byte("BAD REPORTS:\n\n"))
	if err != nil {
		log.L.Warnf("failed to write to file: %s", err)
	}

	for _, report := range compiled {
		if !report.Good {
			_, err = f.Write([]byte(fmt.Sprintf("%v$ %v\n", report.ID, report.Message)))
			if err != nil {
				log.L.Warnf("failed to write to file: %s", err)
				continue
			}
		}
	}

	_, err = f.Write([]byte("\n\n\nGOOD REPORTS:\n\n"))
	if err != nil {
		log.L.Warnf("failed to write to file: %s", err)
	}

	for _, report := range compiled {
		if report.Good {
			_, err = f.Write([]byte(fmt.Sprintf("%v$ %v\n", report.ID, report.Message)))
			if err != nil {
				log.L.Warnf("failed to write to file: %s", err)
				continue
			}
		}
	}

	log.L.Infof("Finished!")
}
