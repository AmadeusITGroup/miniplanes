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

	"github.com/davecgh/go-spew/spew"

	//"k8s.io/client-go/third_party/forked/golang/template"
	airlinesclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airlines"
	"github.com/amadeusitgroup/miniplanes/storage/pkg/gen/models"
	"github.com/amadeusitgroup/miniplanes/ui/cmd/config"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

// AddAirlineHandler handles view to add airline
func AddAirlineHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("AddAirlineHandler")

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

	render(w, r, addAirlineTpl, "add_airline", fullData)

}

// SaveAirlineHandler  handles view to save (if no error) airline
func SaveAirlineHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("SaveAirlineHandler")
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
		"item":          "airline",
	}
	log.Debugf("Request -> %s", spew.Sdump(r))
	statusCode := 201
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Cannot parse request %v: %v", r, err)
		data["err"] = err.Error()
		statusCode = 400
	} else {
		errors := []string{}
		storageURL := fmt.Sprintf("%s:%d", config.StorageHost, config.StoragePort)
		client := airlinesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		params := airlinesclient.NewAddAirlineParams()
		params.Airline = new(models.Airline)
		params.Airline.AirlineID, err = stringToInt64(r.FormValue("AirlineID"))
		if err != nil {
			statusCode = 400
			errors = append(errors, fmt.Sprintf("AirlineID: %v", err))
		}
		params.Airline.Name = r.FormValue("Name")
		params.Airline.Alias = r.FormValue("Alias")
		params.Airline.IATA = r.FormValue("IATA")
		params.Airline.ICAO = r.FormValue("ICAO")
		params.Airline.Callsign = r.FormValue("Callsign")
		params.Airline.Country = r.FormValue("Country")
		params.Airline.Active = r.FormValue("Active")
		log.Debugf("params %s", spew.Sdump(params))
		if len(errors) > 0 {
			data["err"] = strings.Join(errors, ", ")
		} else {
			if _, err := client.AddAirline(params); err != nil {
				data["err"] = err.Error()
			} else {
				log.Infof("Airline added")
			}
		}
	}
	w.WriteHeader(statusCode)
	render(w, r, saveItemTPL, "save_item", data)

}
