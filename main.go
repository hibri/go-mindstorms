package main

import (
	"fmt"
	"time"

	"github.com/hibri/GoEV3/Motor"
	"github.com/hibri/GoEV3/Sensors"
)

func main() {
	sensor := Sensors.FindInfraredSensor(Sensors.InPort2)
	fmt.Printf("Current speed Motor A %d", Motor.CurrentSpeed(Motor.OutPortA))
	Motor.Run(Motor.OutPortA, 40)
	Motor.Run(Motor.OutPortB, 40)

	for {
		value := sensor.ReadProximity()

		if value < 50 {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	Motor.Stop(Motor.OutPortA)
	Motor.Stop(Motor.OutPortB)
}
