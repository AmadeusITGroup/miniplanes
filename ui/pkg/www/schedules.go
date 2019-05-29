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
	"fmt"
	"html/template"
	"net/http"
	"strings"

	schedulesclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/schedules"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/ui/cmd/config"
	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

// AddScheduleHandler handles view to add schedule
func AddScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("AddScheduleHandler")

	push(w, "/static/style.css")
	push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fullData := map[string]interface{}{
		"NavigationBar": template.HTML(navigationBarHTML),
	}
	// pass from global variables
	if !BasicAuth(w, r, username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin protected."`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	render(w, r, addScheduleTpl, "add_schedule", fullData)

}

// SaveScheduleHandler  handles view to save (if no error) schedule
func SaveScheduleHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("SaveScheduleHandler")
	push(w, "/static/style.css")
	push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// pass from global variables
	if !BasicAuth(w, r, username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin protected."`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	data := map[string]interface{}{
		"NavigationBar": template.HTML(navigationBarHTML),
		"item":          "schedule",
	}
	err := r.ParseForm()
	statusCode := 201
	if err != nil {
		log.Errorf("Cannot parse request %v: %v", r, err)
		data["err"] = err.Error()
		statusCode = 400
	} else {
		errors := []string{}
		storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
		client := schedulesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		params := schedulesclient.NewAddScheduleParams()
		params.Schedule = new(models.Schedule)
		log.Debugf("params %+v", params)
		params.Schedule.ArrivalTime = r.FormValue("ArrivalTime")
		params.Schedule.ArriveNextDay, err = stringToBool(r.FormValue("ArriveNextDay"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("ArriveNextDay: %v", err))
		}
		params.Schedule.DaysOperated = r.FormValue("DaysOperated")
		params.Schedule.DepartureTime = r.FormValue("DepartureTime")
		params.Schedule.Destination, err = stringToInt32(r.FormValue("Destination"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Destination: %v", err))
		}
		params.Schedule.FlightNumber = r.FormValue("FlightNumber")
		params.Schedule.OperatingCarrier = r.FormValue("OperatingCarrier")
		params.Schedule.Origin, err = stringToInt32(r.FormValue("Origin"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Origin: %v", err))
		}
		params.Schedule.ScheduleID, err = stringToInt64(r.FormValue("ScheduleID"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("ScheduleID: %v", err))
		}
		if len(errors) > 0 {
			data["err"] = strings.Join(errors, ", ")
		} else {
			if _, err := client.AddSchedule(params); err != nil {
				data["err"] = err.Error()
			} else {
				log.Infof("Schedule added")
			}
		}
	}
	w.WriteHeader(statusCode)
	render(w, r, saveItemTPL, "save_item", data)
}
