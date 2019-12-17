# ds18b20 for raspberry
Simple module to get sensor data of ds18b20 sensor(s) on a Raspberry PI (GPIO w1 pin).
It supports multiple sensors (on bus `w1_bus_master1`) and respects the CRC when reading.

## Connect ds18b20
You'll need to enable onewire interface for this to work.
Add the following to your `/boot/config.txt` to do that:
```
dtoverlay=w1-gpio-pullup,gpiopin=4
```
This configures onewire interface on GPIO4 (connector pin 7), where you should connect your sensor(s) to.

Additionally, the following kernel modules need to be enabled:
```
modprobe wire
modprobe w1-gpio
modprobe w1-therm
```

## Install
go get github.com/marsu-p/ds18b20

## Usage
```
package main

import (
    "fmt"

    "github.com/marsu-p/ds18b20"
)

func main() {
    sensors, err := ds18b20.GetSensors()
    if err != nil {
        panic(err)
    }

    for _, sensor := range sensors {
        t, err := sensor.GetTemperature()
        if err == nil {
            fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor.Name, t)
        }
    }
}
```
