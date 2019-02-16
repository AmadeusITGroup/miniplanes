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
package e2e

import (
	"fmt"

	itinerariesclient "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/client/itineraries"
	livenessclient "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/client/liveness"
	readinessclient "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/gen/client/readiness"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("itineraries-server", func() {
	It("Should have liveness probe", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := livenessclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetLive(nil)
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})

	It("Should have readiness probe", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := readinessclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetReady(readinessclient.NewGetReadyParams())
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should get no itineraries... No To", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := itinerariesclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			params := itinerariesclient.NewGetItinerariesParams()
			boh := "BOH"
			params.From = &boh
			_, err := client.GetItineraries(params)
			if err != nil { // TODO: reflect.equal error
				return nil
			}
			return fmt.Errorf("expected error")
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should get no itineraries... No DepartureDate", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := itinerariesclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			params := itinerariesclient.NewGetItinerariesParams()
			boh := "BOH"
			params.From = &boh
			params.To = &boh
			_, err := client.GetItineraries(params)
			if err != nil { // TODO: reflect.equal error
				return nil
			}
			return fmt.Errorf("expected error")
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should get no itineraries... No ReturnDate", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := itinerariesclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			params := itinerariesclient.NewGetItinerariesParams()
			boh := "BOH"
			params.From = &boh
			params.To = &boh
			bth := "2412"
			params.DepartureDate = &bth
			_, err := client.GetItineraries(params)
			if err != nil { // TODO: reflect.equal error
				return nil
			}
			return fmt.Errorf("expected error")
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should get no itineraries... Can't find airportID for BOH", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := itinerariesclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			params := itinerariesclient.NewGetItinerariesParams()
			boh := "BOH"
			params.From = &boh
			params.To = &boh
			bth := "2412"
			params.DepartureDate = &bth
			params.ReturnDate = &bth
			_, err := client.GetItineraries(params)
			if err != nil { // TODO: reflect.equal error
				return nil
			}
			return fmt.Errorf("expected error")
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should get itineraries...", func() {
		itinerariesServerURL := fmt.Sprintf("%s:%d", ItinerariesServerHost, ItinerariesServerPort)
		client := itinerariesclient.New(httptransport.New(itinerariesServerURL, "", nil), strfmt.Default)
		Eventually(func() error {
			params := itinerariesclient.NewGetItinerariesParams()
			nce := "NCE"
			params.From = &nce
			jfk := "JFK"
			params.To = &jfk
			bth := "2412"
			params.DepartureDate = &bth
			back := "3012"
			params.ReturnDate = &back
			OK, err := client.GetItineraries(params)
			if err != nil { // TODO: reflect.equal error
				return err
			}
			if len(OK.Payload) == 0 {
				return fmt.Errorf("no itineraries found")
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
})
