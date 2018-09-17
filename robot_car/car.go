package main

import (
	"fmt"
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Engine struct {
	motors [4]rpio.Pin
	ena    rpio.Pin
	enb    rpio.Pin

	testMode bool
}

func (e *Engine) init(motor1, motor2, motor3, motor4 int, ena, enb int, test bool) {
	e.motors[0] = rpio.Pin(motor1)
	e.motors[1] = rpio.Pin(motor2)
	e.motors[2] = rpio.Pin(motor3)
	e.motors[3] = rpio.Pin(motor4)

	e.ena = rpio.Pin(ena)
	e.enb = rpio.Pin(enb)
	e.testMode = test

	if !test {
		err := rpio.Open()
		if err != nil {
			log.Panic("can not initialize gpio")
		}
		for _, motor := range e.motors {
			motor.Output()
			motor.Low()
		}

		e.ena.High()
		e.enb.High()

		rpio.StartPwm()
		e.ena.Pwm()
		e.enb.Pwm()
		e.ena.DutyCycle(90, 100)
		e.enb.DutyCycle(90, 100)
		e.ena.Freq(100000)
		e.enb.Freq(100000)
	}
}

func (e *Engine) cleanup() {
	if !e.testMode {
		rpio.StopPwm()
		rpio.Close()
	}
}

func (e *Engine) changeSpeed(speed int) {
	if !e.testMode {
		if speed > 35 && speed <= 100 {
			e.ena.DutyCycle(uint32(speed), 100)
			e.enb.DutyCycle(uint32(speed), 100)
		}
	}
	fmt.Printf("changing speed to [%v]\n", speed)
}

func (e *Engine) move(direction string, interval int) {
	motorsMap := map[string]([2]int){
		"up":    [2]int{0, 2},
		"down":  [2]int{1, 3},
		"left":  [2]int{1, 2},
		"right": [2]int{0, 3},
	}

	if !e.testMode {
		e.ena.High()
		e.enb.High()
		for _, motor := range motorsMap[direction] {
			e.motors[motor].High()
		}
	}
	fmt.Printf("moving car in the [%v] direction\n", direction)
	time.Sleep(time.Millisecond * time.Duration(interval))
}

func (e *Engine) stop() {
	if !e.testMode {
		e.ena.Low()
		e.enb.Low()
		for _, motor := range e.motors {
			motor.Low()
		}
	}
	fmt.Println("stopping the engines")
}

type KunmanCar struct {
	Engine
}

func NewKunmanCar(test bool) *KunmanCar {
	car := new(KunmanCar)
	car.init(19, 16, 21, 26, 13, 20, test)

	return car
}
