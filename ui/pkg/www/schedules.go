package www

import (
	"bytes"
	"net/http"

	storagemodels "github.com/amadeusitgroup/miniapp/storage/pkg/gen/models"
	strfmt "github.com/go-openapi/strfmt"
	//"k8s.io/client-go/third_party/forked/golang/template"
)

var (
	origin           = int64(532)
	destination      = int64(1382)
	flightNumber     = "9W4777"
	operatingCarrier = "Air France"
	departure        = strfmt.DateTime{}
	arrival          = strfmt.DateTime{}
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
			Departure:        &departure,
			Arrival:          &arrival,
		},
	}
	w.WriteHeader(200)
	buf := new(bytes.Buffer)
	schedulesViewTpl.Execute(buf, schedules)
	w.Write(buf.Bytes())
}
