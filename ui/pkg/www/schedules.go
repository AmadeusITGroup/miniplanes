/*
Copyright 2018 Amadeus SaS All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package www

import (
	"bytes"
	"net/http"

	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	//"k8s.io/client-go/third_party/forked/golang/template"
)

var (
	origin           = int32(532)
	destination      = int32(1382)
	flightNumber     = "9W4777"
	operatingCarrier = "Air France"
	departure        = "1011"
	arrival          = "1823"
	daysOperated     = "1234567"
)

// SearchSchedules search schedules
func SearchSchedules(w http.ResponseWriter, r *http.Request) {
	push(w, "/static/style.css")
	push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		// log.Error
		return
	}
	// Here we should get the schedules...

	//532,1382,9W4777,Air France,1234567,0600,0905
	// origin destination flightNumber operatingCarrier daysOperated departure arrival
	schedules := []storagemodels.Schedule{
		storagemodels.Schedule{
			Origin:           &destination,
			Destination:      &origin,
			FlightNumber:     &flightNumber,
			OperatingCarrier: &operatingCarrier,
			DaysOperated:     &daysOperated,
			DepartureTime:    &departure,
			ArrivalTime:      &arrival,
		},
	}
	w.WriteHeader(200)
	buf := new(bytes.Buffer)
	schedulesViewTpl.Execute(buf, schedules)
	w.Write(buf.Bytes())
}
