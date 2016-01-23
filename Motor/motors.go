// Provides APIs for interacting with EV3's motors.
package Motor

import (
	"github.com/ldmberman/GoEV3/utilities"
	"log"
	"os"
	"path"
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
	portFD           = "address"
	regulationModeFD = "speed_regulation"
	speedGetterFD    = "speed"
	speedSetterFD    = "speed_sp"
	powerGetterFD    = "duty_cycle"
	powerSetterFD    = "duty_cycle_sp"
	runFD            = "command"
	stopModeFD       = "stop_command"
	positionFD       = "position"
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

	log.Fatal("No motor is connected to port ", port )
	return ""
}

// Runs the motor at the given port.
// The meaning of `speed` parameter depends on whether the regulation mode is turned on or off.
//
// When the regulation mode is off (by default) `speed` ranges from -100 to 100 and
// it's absolute value indicates the percent of motor's power usage. It can be roughly interpreted as
// a motor speed, but deepending on the environment, the actual speed of the motor
// may be lower than the target speed.
//
// When the regulation mode is on (has to be enabled by EnableRegulationMode function) the motor
// driver attempts to keep the motor speed at the `speed` value you've specified
// which ranges from about -1000 to 1000. The actual range depends on the type of the motor - see ev3dev docs.
//
// Negative values indicate reverse motion regardless of the regulation mode.
func Run(port OutPort, speed int16) {
	folder := findFolder(port)
	regulationMode := utilities.ReadStringValue(folder, regulationModeFD)

	switch regulationMode {
	case "on":
		utilities.WriteIntValue(folder, speedSetterFD, int64(speed))
		utilities.WriteStringValue(folder, runFD, "run-forever")
	case "off":
		if speed > 100 || speed < -100 {
			log.Fatal("The speed must be in range [-100, 100]")
		}
		utilities.WriteIntValue(folder, powerSetterFD, int64(speed))
		utilities.WriteStringValue(folder, runFD, "run-forever")
	}
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

// Enables regulation mode, causing the motor at the given port to compensate
// for any resistance and maintain its target speed.
func EnableRegulationMode(port OutPort) {
	utilities.WriteStringValue(findFolder(port), regulationModeFD, "on")
}

// Disables regulation mode. Regulation mode is off by default.
func DisableRegulationMode(port OutPort) {
	utilities.WriteStringValue(findFolder(port), regulationModeFD, "off")
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
