package Sensors

import (
	"fmt"
	"github.com/ldmberman/GoEV3/utilities"
)

// Gyro sensor type.
type GyroSensor struct {
	port InPort
}

// Provides access to a gyro sensor at the given port.
func FindGyroSensor(port InPort) *GyroSensor {
	snr := findSensor(port, TypeGyro)

	s := new(GyroSensor)
	s.port = port

	path := fmt.Sprintf("%s/%s", baseSensorPath, snr)
	utilities.WriteStringValue(path, "mode", "GYRO-G&A")

	return s
}

// Reads the angle of degrees.
func (self *GyroSensor) ReadAngle() int16 {
	snr := findSensor(self.port, TypeGyro)

	path := fmt.Sprintf("%s/%s", baseSensorPath, snr)
	value := utilities.ReadInt16Value(path, "value0")

	return value
}

// Reads the rotational speed in range [-440, 440].
func (self *GyroSensor) ReadRotationalSpeed() int16 {
	snr := findSensor(self.port, TypeGyro)

	path := fmt.Sprintf("%s/%s", baseSensorPath, snr)
	value := utilities.ReadInt16Value(path, "value1")

	return value
}
