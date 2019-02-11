/*

MIT License

Copyright (c) 2019 Amadeus s.a.s.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

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
