// Copyright 2017 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 {the "License"};
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/davecheney/gpio"
	"github.com/mendersoftware/thermostat/hwio"
	"github.com/mendersoftware/thermostat/server"
	"github.com/mendersoftware/thermostat/store"
)

var data store.MemStore

func main() {
	// init storeage
	data = store.InitStore()

	distC := make(chan bool)

	// store sensor readings
	readings := make(chan hwio.HumiTemp)
	go func() {
		for {
			humiTemp, ok := <-readings
			if ok {
				data.Store("temp", humiTemp.Temp)
				data.Store("humi", humiTemp.Humi)
			} else {
				return
			}
		}
	}()

	go readDHT(readings)

	buzzerPin, err := hwio.PinToGpioExport("P8_8")
	if err != nil {
		fmt.Println(err)
	}
	buzzer, err := gpio.OpenPin(int(buzzerPin), gpio.ModeOutput)
	if err != nil {
		fmt.Printf("Error initializing buzzer: %s\n", err)
		buzzer = nil
	} else {
		// set the high state so that the buzzer will be silent
		buzzer.Set()
	}

	distancePin, err := hwio.PinToGpioExport("P8_10")
	if err != nil {
		fmt.Println(err)
	}

	dsitSensor, err := gpio.OpenPin(int(distancePin), gpio.ModeInput)
	if err != nil {
		fmt.Printf("Error initializing distance sensor: %s\n", err)
	} else {
		err = dsitSensor.BeginWatch(gpio.EdgeBoth, func() {
			// check if we have a buzzer to make soem sound
			if buzzer != nil {
				if dsitSensor.Get() {
					buzzer.Clear()
				} else {
					buzzer.Set()
				}
			}

			// do not block while sending message
			select {
			case distC <- dsitSensor.Get():
				// do nothing
			default:
				// do nothing here as well
			}
		})
		if err != nil {
			fmt.Printf("Unable to watch pin: %s\n", err.Error())
		}
	}

	if err := startServer(data, distC); err != nil {
		fmt.Printf("Unable to start server: %s\n", err)
	}
}

func readDHT(dataC chan hwio.HumiTemp) {
	for {
		tempHumi, err := hwio.ReadHumiTempWithRetry("P8_11")
		if err != nil {
			fmt.Println(err)
		}

		// print temperature and humidity
		fmt.Printf("Temperature = %v*C, Humidity = %v%%\n",
			tempHumi.Temp, tempHumi.Humi)

		dataC <- tempHumi

		// wait a bit before reading again
		time.Sleep(time.Second * 5)
	}
}

func startServer(data store.MemStore, statusC <-chan bool) error {
	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "80"
	}

	handler := &server.RegexpHandler{}

	handler.HandleFunc(regexp.MustCompile("^/history"), showHistory)
	handler.HandleFunc(regexp.MustCompile("^/export/[A-z]"), doExport)
	handler.HandleFunc(regexp.MustCompile("^/distance"), showDistance)
	handler.HandleFunc(regexp.MustCompile("^/ws"),
		func(w http.ResponseWriter, r *http.Request) {
			handleWS(w, r, statusC)
		})
	handler.HandleFunc(regexp.MustCompile("^/$"), showWeather)

	fmt.Printf("Server is running on port %v\n", port)

	return http.ListenAndServe(":"+port, handler)
}
