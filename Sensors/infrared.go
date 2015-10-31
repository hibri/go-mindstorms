package Sensors

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ldmberman/GoEV3/utilities"
)

const (
	RedUp    Button = 1
	RedDown         = 2
	BlueUp          = 3
	BlueDown        = 4

	Channel1 Channel = 0
	Channel2         = 1
	Channel3         = 2
	Channel4         = 3
)

var (
	REMOTE_POLLING_INTERVAL = 500 // milliseconds
)

type (
	// Infrared sensor type.
	InfraredSensor struct {
		port InPort
	}

	RemoteSignal struct {
		Name  string
		Value uint64
	}
	Button  uint64
	Channel uint64
)

// Provides access to an infrared sensor at the given port.
func FindInfraredSensor(port InPort) *InfraredSensor {
	findSensor(port, TypeInfrared)

	s := new(InfraredSensor)
	s.port = port

	return s
}

// Reads the proximity value (in range 0 - 100) reported by the infrared sensor. A value of 100 corresponds to a range of approximately 70 cm.
func (self *InfraredSensor) ReadProximity() uint8 {
	snr := findSensor(self.port, TypeInfrared)

	path := fmt.Sprintf("%s/%s", baseSensorPath, snr)
	utilities.WriteStringValue(path, "mode", "IR-PROX")
	value := utilities.ReadUInt8Value(path, "value0")

	return value
}

// Blocks until the infrared sensor detects a nearby object.
func (self *InfraredSensor) WaitForProximity() {
	for {
		p1 := self.ReadProximity()
		time.Sleep(time.Millisecond * 100)
		p2 := self.ReadProximity()

		if p1 < 20 && p2 < 20 {
			return
		}
	}
}

// Turns on the remote control mode.
func (self *InfraredSensor) RemoteModeOn() {
	snr := findSensor(self.port, TypeInfrared)
	path := fmt.Sprintf("%s/%s", baseSensorPath, snr)
	utilities.WriteStringValue(path, "mode", "IR-REMOTE")
}

// Registers a callback to be triggered when a remote button is pressed. The listening
// can be stopped by sending any boolean value to a `stop` channel.
func (self *InfraredSensor) OnRemotePressed(stop <-chan bool, fn func(c Channel, b Button)) {
	self.RemoteModeOn()
	s := make(chan RemoteSignal, 50)

	go func() {
		pressed := map[uint64]bool{}
		for {
			select {
			case <-stop:
				return
			case signal := <-s:
				c := parseChannel(signal.Name)

				if signal.Value == 0 {
					for b := RedUp; b <= BlueDown; b++ {
						pressed[buttonID(c, b)] = false
					}
					continue
				}
				k := buttonID(c, Button(signal.Value))
				if v, ok := pressed[k]; ok && v {
					continue
				}
				pressed[k] = true
				fn(c, Button(signal.Value))
			}
		}
	}()
	self.pollRemote(s, stop)
}

// Registers a callback to be triggered when a remote button is released. The listening
// can be stopped by sending any boolean value to a `stop` channel.
func (self *InfraredSensor) OnRemoteReleased(stop <-chan bool, fn func(c Channel, b Button)) {
	self.RemoteModeOn()
	s := make(chan RemoteSignal, 50)

	go func() {
		pressed := map[uint64]bool{}
		for {
			select {
			case <-stop:
				return
			case signal := <-s:
				c := parseChannel(signal.Name)

				if signal.Value != 0 {
					pressed[buttonID(c, Button(signal.Value))] = true
					continue
				}
				for b := RedUp; b <= BlueDown; b++ {
					if v, ok := pressed[buttonID(c, b)]; ok && v {
						fn(c, b)
						pressed[buttonID(c, b)] = false
					}
				}
			}
		}
	}()
	self.pollRemote(s, stop)
}

func parseChannel(name string) Channel {
	var c Channel
	switch name {
	case "value0":
		c = Channel1
	case "value1":
		c = Channel2
	case "value2":
		c = Channel3
	case "value3":
		c = Channel4
	default:
		log.Fatal("Invalid channel")
	}
	return c
}

func buttonID(c Channel, b Button) uint64 {
	return uint64(c)*10 + uint64(b)
}

func (self *InfraredSensor) pollRemote(s chan<- RemoteSignal, stop <-chan bool) {
	snr := findSensor(self.port, TypeInfrared)
	for i := 0; i < 4; i++ {
		name := fmt.Sprintf("value%d", i)
		p := fmt.Sprintf("%s/%s/%s", baseSensorPath, snr, name)
		go func() {
			f, err := os.Open(p)
			defer f.Close()
			if err != nil {
				log.Fatal(err)
			}
			for {
				select {
				case <-stop:
					return
				default:
				}

				data, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal(err)
				}
				_, err = f.Seek(0, 0)
				if err != nil {
					log.Fatal(err)
				}
				b, err := strconv.ParseUint(strings.Trim(string(data), " \n"), 10, 16)
				if err != nil {
					log.Fatal(err)
				}
				s <- RemoteSignal{name, b}
				time.Sleep(time.Millisecond * time.Duration(REMOTE_POLLING_INTERVAL))
			}
		}()
	}
}
