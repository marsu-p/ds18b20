package ds18b20

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var busPath = "/sys/bus/w1/devices/"
var busMasterPath = busPath + "w1_bus_master1/"
var busMaster = busMasterPath + "w1_master_slaves"
var sensorFile = "w1_slave"

// Sensor type
type Sensor struct {
	Name string
}

// GetSensors get known sensors
func GetSensors() (sensors []*Sensor, err error) {
	data, err := ioutil.ReadFile(busMaster)
	if err != nil {
		return nil, err
	}

	for _, sensorName := range strings.Split(string(data), "\n") {
		if sensorName == "" {
			continue
		}
		sensors = append(sensors, &Sensor{Name: sensorName})
	}

	return sensors, nil
}

// GetTemperature get the temperature for a sensor
//
// the sensor device's content looks like this:
// 33 00 4b 46 ff ff 02 10 f4 : crc=f4 YES
// 33 00 4b 46 ff ff 02 10 f4 t=25625
func (s *Sensor) GetTemperature() (temp float64, err error) {
	data, err := ioutil.ReadFile(busPath + s.Name + "/w1_slave")
	if err != nil {
		return 0.0, fmt.Errorf("failed to read sensor %q: %v", s.Name, err)
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) != 3 {
		return 0.0, fmt.Errorf("can not parse temperature, wrong file format")
	}

	// 49 01 4b 46 7f ff 07 10 f6 : crc=f6 YES
	regex := regexp.MustCompile(`(?:[0-9a-f]{2} ){9}: crc=(?:[0-9a-f]{2} )YES`)
	if !regex.MatchString(lines[0]) {
		return 0.0, fmt.Errorf("could not verify crc, regex did not match")
	}

	// 33 00 4b 46 ff ff 02 10 f4 t=25625
	regex = regexp.MustCompile(`(?:[0-9a-f]{2} ){9}t=(?P<temperature>[0-9]+)`)

	match := regex.FindStringSubmatch(lines[1])
	if match == nil {
		return 0.0, fmt.Errorf("could not extract temperature, regex did not match")
	}
	temperature, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return 0.0, fmt.Errorf("could not extract temperature, parse to float failed")
	}
	temp = float64(temperature) / 1000

	return temp, nil
}
