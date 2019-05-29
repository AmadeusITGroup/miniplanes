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

	//"k8s.io/client-go/third_party/forked/golang/template"

	airportsclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airports"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/ui/cmd/config"
	"github.com/davecgh/go-spew/spew"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

// AddAirportHandler handles view to add Airport
func AddAirportHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("AddAirportHandler")

	push(w, "/static/style.css")
	push(w, "/static/navigation_bar.css")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := map[string]interface{}{
		"NavigationBar": template.HTML(navigationBarHTML),
	}
	// pass from global variables
	if !BasicAuth(w, r, username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Admin protected."`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	render(w, r, addAirportTpl, "add_airport", data)

}

// SaveAirportHandler  handles view to save (if no error) Airport
func SaveAirportHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("SaveAirportHandler")
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
		"item":          "airport",
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
		client := airportsclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		params := airportsclient.NewAddAirportParams()
		params.Airport = new(models.Airport)
		params.Airport.AirportID, err = stringToInt32(r.FormValue("AirportID"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("AirportID: %v", err))
		}
		params.Airport.Altitude, err = stringToFloat64(r.FormValue("Altitude"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Altitude: %v", err))
		}
		params.Airport.City = r.FormValue("City")
		params.Airport.Country = r.FormValue("Country")
		params.Airport.DST = r.FormValue("DST")
		params.Airport.Latitude, err = stringToFloat64(r.FormValue("Latitude"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Latitude: %v", err))
		}
		params.Airport.Longitude, err = stringToFloat64(r.FormValue("Longitude"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Longitude: %v", err))
		}
		params.Airport.Name = r.FormValue("Name")
		params.Airport.TZ = r.FormValue("TZ")
		params.Airport.Timezone, err = stringToInt64(r.FormValue("Timezone"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("Timezone: %v", err))
		}
		log.Debugf("params %s", spew.Sdump(params))
		if len(errors) > 0 {
			data["err"] = strings.Join(errors, ", ")
		} else {
			if _, err := client.AddAirport(params); err != nil {
				data["err"] = err.Error()
			} else {
				log.Infof("Airport added")
			}
		}
	}

	w.WriteHeader(statusCode)
	render(w, r, saveItemTPL, "save_item", data)

}
