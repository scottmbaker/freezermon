package ds18b20

// DS18B20 collector
// Scott Baker
//
// Collects DS18B20 measurements by reading the files in /sys/bus/w1/devices
//
// This is a hasty port of the python library I wrote to do this a long, long time ago.
// We don't need to do GPIO ourselves -- the 1wire driver does it for us. We just need
// to read our sensor from /sys/bus/w1/devices.
//
// We assume there is only one sensor. If there is more than one, this it will pick the
// first one, for some definition of first.

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	DEVICE_DIR = "/sys/bus/w1/devices"
)

type DS18B20 struct {
	Devices []string
	verbose bool
}

func (ds *DS18B20) FindDevices() error {
	ds.Devices = []string{}

	files, err := os.ReadDir(DEVICE_DIR)
	if err != nil {
		return fmt.Errorf("error reading device directory %s: %v", DEVICE_DIR, err)
	}

	for _, file := range files {
		if ds.verbose {
			log.Printf("checking file: %s\n", file.Name())
		}
		if file.Name()[:3] == "10-" || file.Name()[:3] == "28-" {
			if ds.verbose {
				log.Printf("Found device: %s\n", file.Name())
			}
			ds.Devices = append(ds.Devices, file.Name())
		}
	}
	return nil
}

func (ds *DS18B20) MeasureDevice(name string) (float64, error) {
	deviceFileName := fmt.Sprintf("%s/%s/w1_slave", DEVICE_DIR, name)
	if _, err := os.Stat(deviceFileName); os.IsNotExist(err) {
		return 0.0, fmt.Errorf("device file %s does not exist", deviceFileName)
	}
	fileContents, err := os.ReadFile(deviceFileName)
	if err != nil {
		return 0.0, fmt.Errorf("error reading device file %s: %v", deviceFileName, err)
	}
	lines := strings.Split(string(fileContents), "\n")
	if len(lines) < 2 {
		return 0.0, fmt.Errorf("unexpected format in device file %s", deviceFileName)
	}

	if !strings.HasSuffix(lines[0], "YES") {
		return 0.0, fmt.Errorf("first line does not end with yes in %s: %s", deviceFileName, lines[0])
	}

	parts := strings.Split(lines[1], "t=")
	if len(parts) != 2 {
		return 0.0, fmt.Errorf("unexpected format in second line of %s: %s", deviceFileName, lines[1])
	}

	tempC, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing temperature in %s: %v", deviceFileName, err)
	}

	return tempC / 1000.0, nil
}

func (ds *DS18B20) MeasureFirstDevice() (float64, error) {
	if len(ds.Devices) == 0 {
		return 0.0, fmt.Errorf("no devices found")
	}
	return ds.MeasureDevice(ds.Devices[0])
}

func (ds *DS18B20) GetDeviceCount() int {
	return len(ds.Devices)
}

func NewDS18B20(verbose bool) (*DS18B20, error) {
	ds := &DS18B20{verbose: verbose}
	err := ds.FindDevices()
	if err != nil {
		return nil, fmt.Errorf("error finding devices: %v", err)
	}
	return ds, nil
}
