package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hibri/GoEV3/Motor"
	"github.com/hibri/GoEV3/Sensors"
)

func main() {
	sensor := Sensors.FindInfraredSensor(Sensors.InPort2)
	colorSensor := Sensors.FindColorSensor(Sensors.InPort1)

	fmt.Printf("Current speed Motor A %d", Motor.CurrentSpeed(Motor.OutPortA))
	fmt.Printf("Current Sensor reading %d", sensor.ReadProximity())

	Motor.Run(Motor.OutPortA, 150)
	Motor.Run(Motor.OutPortB, 150)

	for {

		color := colorSensor.ReadColor()
		log.Print(color)
		if color == Sensors.Red {
			Motor.Stop(Motor.OutPortA)
			Motor.Stop(Motor.OutPortB)
		}
		if color == Sensors.White {
			Motor.Run(Motor.OutPortA, -100)
			Motor.Run(Motor.OutPortB, -100)
		}
		value := sensor.ReadProximity()

		if value < 50 {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
}
