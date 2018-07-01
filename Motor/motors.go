// Provides APIs for interacting with EV3's motors.
package Motor

import (
	"log"
	"os"
	"path"

	"github.com/hibri/GoEV3/utilities"
)

// Constants for output ports.
type OutPort string

const (
	OutPortA OutPort = "A"
	OutPortB         = "B"
	OutPortC         = "C"
	OutPortD         = "D"
)

// Names of files which constitute the low-level motor API
const (
	rootMotorPath = "/sys/class/tacho-motor"
	// File descriptors for getting/setting parameters
	portFD        = "address"
	speedGetterFD = "speed"
	speedSetterFD = "speed_sp"
	powerGetterFD = "duty_cycle"
	powerSetterFD = "duty_cycle_sp"
	runFD         = "command"
	stopModeFD    = "stop_command"
	positionFD    = "position"
)

func findFolder(port OutPort) string {
	if _, err := os.Stat(rootMotorPath); os.IsNotExist(err) {
		log.Fatal("There are no motors connected")
	}

	rootMotorFolder, _ := os.Open(rootMotorPath)
	motorFolders, _ := rootMotorFolder.Readdir(-1)
	if len(motorFolders) == 0 {
		log.Fatal("There are no motors connected")
	}

	for _, folderInfo := range motorFolders {
		folder := folderInfo.Name()
		motorPort := utilities.ReadStringValue(path.Join(rootMotorPath, folder), portFD)
		if motorPort == "out"+string(port) {
			return path.Join(rootMotorPath, folder)
		}
	}

	log.Fatal("No motor is connected to port ", port)
	return ""
}

// Runs the motor at the given port, untill stop is called
//
// Negative values indicate reverse motion
func Run(port OutPort, speed int16) {
	folder := findFolder(port)

	utilities.WriteIntValue(folder, speedSetterFD, int64(speed))
	utilities.WriteStringValue(folder, runFD, "run-forever")

}

// Stops the motor at the given port.
func Stop(port OutPort) {
	utilities.WriteStringValue(findFolder(port), runFD, "stop")
}

// Reads the operating speed of the motor at the given port.
func CurrentSpeed(port OutPort) int16 {
	return utilities.ReadInt16Value(findFolder(port), speedGetterFD)
}

// Reads the operating power of the motor at the given port.
func CurrentPower(port OutPort) int16 {
	return utilities.ReadInt16Value(findFolder(port), powerGetterFD)
}

// Enables brake mode, causing the motor at the given port to brake to stops.
func EnableBrakeMode(port OutPort) {
	utilities.WriteStringValue(findFolder(port), stopModeFD, "brake")
}

// Disables brake mode, causing the motor at the given port to coast to stops. Brake mode is off by default.
func DisableBrakeMode(port OutPort) {
	utilities.WriteStringValue(findFolder(port), stopModeFD, "coast")
}

// Reads the position of the motor at the given port.
func CurrentPosition(port OutPort) int32 {
	return utilities.ReadInt32Value(findFolder(port), positionFD)
}

// Set the position of the motor at the given port.
func InitializePosition(port OutPort, value int32) {
	utilities.WriteIntValue(findFolder(port), positionFD, int64(value))
}
