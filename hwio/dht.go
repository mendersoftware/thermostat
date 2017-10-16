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

/*
#include "bbb_dht_read.h"
*/
import "C"
import (
	"errors"
	"time"
)

type HumiTemp struct {
	Humi float32
	Temp float32
}

const DHT11 = 11

func ReadHumiTempWithRetry(name string) (HumiTemp, error) {
	pin, err := PinToGpio(name)
	if err != nil {
		return HumiTemp{}, err
	}

	var avgTemp, avgHumi float32
	success := 0

	for i := 0; i < 5; i++ {
		temp := C.float(0)
		humi := C.float(0)
		ret := C.bbb_dht_read(C.int(DHT11), C.int(pin.base),
			C.int(pin.number), &humi, &temp)

		if int(ret) == 0 && temp > 1 && humi > 1 {
			success++
			avgTemp = avgTemp + float32(temp)
			avgHumi = avgHumi + float32(humi)
		}
		time.Sleep(time.Microsecond * 100)
	}

	if success == 0 {
		return HumiTemp{}, errors.New("error reading DHT11")
	}
	return HumiTemp{avgHumi / float32(success), avgTemp / float32(success)}, nil
}
