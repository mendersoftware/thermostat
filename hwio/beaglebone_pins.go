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

package hwio

import "errors"

const High = true
const Low = false

type GpioPin struct {
	base   int
	number int
}

func PinToGpio(name string) (GpioPin, error) {
	gpioPin, ok := BeagleboneGPIOPinMap[name]
	if !ok {
		return GpioPin{}, errors.New("invalid pin")
	}
	return gpioPin, nil
}

func PinToGpioExport(name string) (uint, error) {
	gpioPin, ok := BeagleboneGPIOPinMap[name]
	if !ok {
		return 0, errors.New("invalid pin")
	}
	// To get to the GPIO number that should be exported from the kernel,
	// we must add 32 to each GPIO from gpio1, 64 for gpio2, and 96 for gpio3.

	return uint(gpioPin.base*32 + gpioPin.number), nil
}

type PinMap map[string]GpioPin

// Define mapping of pin names to GPIO base and number.
// Adapted from Adafruit_BBIO code Beaglebone Black system reference.
var BeagleboneGPIOPinMap = PinMap{
	"P9_11":      {0, 30},
	"P9_12":      {1, 28},
	"P9_13":      {0, 31},
	"P9_14":      {1, 18},
	"P9_15":      {1, 16},
	"P9_16":      {1, 19},
	"P9_17":      {0, 5},
	"P9_18":      {0, 4},
	"P9_19":      {0, 13},
	"P9_20":      {0, 12},
	"P9_21":      {0, 3},
	"P9_22":      {0, 2},
	"P9_23":      {1, 17},
	"P9_24":      {0, 15},
	"P9_25":      {3, 21},
	"P9_26":      {0, 14},
	"P9_27":      {3, 19},
	"P9_28":      {3, 17},
	"P9_29":      {3, 15},
	"P9_30":      {3, 16},
	"P9_31":      {3, 14},
	"P9_41":      {0, 20},
	"P9_42":      {0, 7},
	"UART4_RXD":  {0, 30},
	"UART4_TXD":  {0, 31},
	"EHRPWM1A":   {1, 18},
	"EHRPWM1B":   {1, 19},
	"I2C1_SCL":   {0, 5},
	"I2C1_SDA":   {0, 4},
	"I2C2_SCL":   {0, 13},
	"I2C2_SDA":   {0, 12},
	"UART2_TXD":  {0, 3},
	"UART2_RXD":  {0, 2},
	"UART1_TXD":  {0, 15},
	"UART1_RXD":  {0, 14},
	"SPI1_CS0":   {3, 17},
	"SPI1_D0":    {3, 15},
	"SPI1_D1":    {3, 16},
	"SPI1_SCLK":  {3, 14},
	"CLKOUT2":    {0, 20},
	"30":         {0, 30},
	"60":         {1, 28},
	"31":         {0, 31},
	"50":         {1, 18},
	"48":         {1, 16},
	"51":         {1, 19},
	"5":          {0, 5},
	"4":          {0, 4},
	"13":         {0, 13},
	"12":         {0, 12},
	"3":          {0, 3},
	"2":          {0, 2},
	"49":         {1, 17},
	"15":         {0, 15},
	"117":        {3, 21},
	"14":         {0, 14},
	"115":        {3, 19},
	"113":        {3, 17},
	"111":        {3, 15},
	"112":        {3, 16},
	"110":        {3, 14},
	"20":         {0, 20},
	"7":          {0, 7},
	"P8_3":       {1, 6},
	"P8_4":       {1, 7},
	"P8_5":       {1, 2},
	"P8_6":       {1, 3},
	"P8_7":       {2, 2},
	"P8_8":       {2, 3},
	"P8_9":       {2, 5},
	"P8_10":      {2, 4},
	"P8_11":      {1, 13},
	"P8_12":      {1, 12},
	"P8_13":      {0, 23},
	"P8_14":      {0, 26},
	"P8_15":      {1, 15},
	"P8_16":      {1, 14},
	"P8_17":      {0, 27},
	"P8_18":      {2, 1},
	"P8_19":      {0, 22},
	"P8_20":      {1, 31},
	"P8_21":      {1, 30},
	"P8_22":      {1, 5},
	"P8_23":      {1, 4},
	"P8_24":      {1, 1},
	"P8_25":      {1, 0},
	"P8_26":      {1, 29},
	"P8_27":      {2, 22},
	"P8_28":      {2, 24},
	"P8_29":      {2, 23},
	"P8_30":      {2, 25},
	"P8_31":      {0, 10},
	"P8_32":      {0, 11},
	"P8_33":      {0, 9},
	"P8_34":      {2, 17},
	"P8_35":      {0, 8},
	"P8_36":      {2, 16},
	"P8_37":      {2, 14},
	"P8_38":      {2, 15},
	"P8_39":      {2, 12},
	"P8_40":      {2, 13},
	"P8_41":      {2, 10},
	"P8_42":      {2, 11},
	"P8_43":      {2, 8},
	"P8_44":      {2, 9},
	"P8_45":      {2, 6},
	"P8_46":      {2, 7},
	"TIMER4":     {2, 2},
	"TIMER7":     {2, 3},
	"TIMER5":     {2, 5},
	"TIMER6":     {2, 4},
	"EHRPWM2B":   {0, 23},
	"EHRPWM2A":   {0, 22},
	"UART5_CTSN": {0, 10},
	"UART5_RTSN": {0, 11},
	"UART4_RTSN": {0, 9},
	"UART3_RTSN": {2, 17},
	"UART4_CTSN": {0, 8},
	"UART3_CTSN": {2, 16},
	"UART5_TXD":  {2, 14},
	"UART5_RXD":  {2, 15},
	"38":         {1, 6},
	"39":         {1, 7},
	"34":         {1, 2},
	"35":         {1, 3},
	"66":         {2, 2},
	"67":         {2, 3},
	"69":         {2, 5},
	"68":         {2, 4},
	"45":         {1, 13},
	"44":         {1, 12},
	"23":         {0, 23},
	"26":         {0, 26},
	"47":         {1, 15},
	"46":         {1, 14},
	"27":         {0, 27},
	"65":         {2, 1},
	"22":         {0, 22},
	"63":         {1, 31},
	"62":         {1, 30},
	"37":         {1, 5},
	"36":         {1, 4},
	"33":         {1, 1},
	"32":         {1, 0},
	"61":         {1, 29},
	"86":         {2, 22},
	"88":         {2, 24},
	"87":         {2, 23},
	"89":         {2, 25},
	"10":         {0, 10},
	"11":         {0, 11},
	"9":          {0, 9},
	"81":         {2, 17},
	"8":          {0, 8},
	"80":         {2, 16},
	"78":         {2, 14},
	"79":         {2, 15},
	"76":         {2, 12},
	"77":         {2, 13},
	"74":         {2, 10},
	"75":         {2, 11},
	"72":         {2, 8},
	"73":         {2, 9},
	"70":         {2, 6},
	"71":         {2, 7},
}
