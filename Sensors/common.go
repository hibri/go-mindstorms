// Provides APIs for interacting with EV3's sensors.
package Sensors

import (
	"fmt"
	"github.com/ldmberman/GoEV3/utilities"
	"io/ioutil"
	"log"
	"strings"
)

// Constants for input ports.
type InPort string

const (
	baseSensorPath = "/sys/class/lego-sensor"

	InPort1 InPort = "in1"
	InPort2        = "in2"
	InPort3        = "in3"
	InPort4        = "in4"
)

// Constants for sensor types.
type Type string

const (
	TypeTouch      Type = "lego-ev3-touch"
	TypeColor           = "lego-ev3-color"
	TypeUltrasonic      = "lego-ev3-us"
	TypeInfrared        = "lego-ev3-ir"
	TypeGyro            = "lego-ev3-gyro"
)

func (self Type) String() string {
	switch self {
	case TypeTouch:
		return "touch"
	case TypeColor:
		return "color"
	case TypeUltrasonic:
		return "ultrasonic"
	case TypeInfrared:
		return "infrared"
	case TypeGyro:
		return "gyro"
	default:
		return "unknown"
	}
}

func findSensor(port InPort, t Type) string {
	sensors, _ := ioutil.ReadDir(baseSensorPath)

	for _, item := range sensors {
		if strings.HasPrefix(item.Name(), "sensor") {
			sensorPath := fmt.Sprintf("%s/%s", baseSensorPath, item.Name())
			portr := utilities.ReadStringValue(sensorPath, "port_name")

			if InPort(portr) == port {
				typer := utilities.ReadStringValue(sensorPath, "driver_name")

				if Type(typer) == t {
					return item.Name()
				}
			}
		}
	}

	log.Fatalf("Could not find %v sensor on port %v\n", t, port)

	return ""
}
