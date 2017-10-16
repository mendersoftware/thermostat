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
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func parseTemplate(name string, data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	template, err := template.ParseFiles(name)
	if err != nil {
		return nil, err
	}
	err = template.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func showWeather(w http.ResponseWriter, r *http.Request) {
	fmt.Println("weather request")
	w.Header().Set("Content-Type", "text/html")

	template, err := parseTemplate("/var/www/weather.html",
		map[string]interface{}{
			"Temperature": data.GetLast("temp").Value,
			"Humidity":    data.GetLast("humi").Value,
			"LastReading": data.GetLast("temp").Date.Format(time.UnixDate),
		})
	if err != nil {
		fmt.Printf("error parsing weather template: %v", err.Error())
		return
	}

	// cookie := http.Cookie{Name: "thermostat",
	// 	Value: "ELC", Expires: time.Now().Add(2 * time.Hour)}
	// http.SetCookie(w, &cookie)

	fmt.Fprintf(w, string(template))
}

func showHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("weather request")
	w.Header().Set("Content-Type", "text/html")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	template, err := parseTemplate("/var/www/history.html",
		map[string]interface{}{
			"Temp": data.GetAll("temp"),
			"Humi": data.GetAll("humi"),
		})
	if err != nil {
		fmt.Printf("error parsing weather template: %v", err.Error())
		return
	}
	fmt.Fprintf(w, string(template))
}

func doExport(w http.ResponseWriter, r *http.Request) {
	fmt.Println("landing page")
	w.Header().Set("Content-Type", "application/json")

	toExport := strings.Split(r.URL.EscapedPath(), "/")

	outputJson, err := json.Marshal(data.GetAll(toExport[len(toExport)-1]))
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "{}")
		return
	}
	fmt.Fprintf(w, string(outputJson))
}
