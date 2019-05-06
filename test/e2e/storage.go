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

	airlinesclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airlines"
	airportsclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/airports"
	livenessclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/liveness"
	readinessclient "github.com/amadeusitgroup/miniplanes/storage/pkg/gen/client/readiness"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("storage", func() {
	It("Should have liveness probe", func() {
		storageURL := fmt.Sprintf("%s:%d", StorageHost, StoragePort)
		client := livenessclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetLive(nil)
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})

	It("Should have readiness probe", func() {
		storageURL := fmt.Sprintf("%s:%d", StorageHost, StoragePort)
		client := readinessclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetReady(readinessclient.NewGetReadyParams())
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})

	It("Should list airports", func() {
		storageURL := fmt.Sprintf("%s:%d", StorageHost, StoragePort)
		client := airportsclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		Eventually(func() error {
			OK, err := client.GetAirports(airportsclient.NewGetAirportsParams())
			if err != nil {
				return err
			}
			if len(OK.Payload) == 0 {
				return fmt.Errorf("no airports found")
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	It("Should list airlines", func() {
		storageURL := fmt.Sprintf("%s:%d", StorageHost, StoragePort)
		client := airlinesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
		Eventually(func() error {
			OK, err := client.GetAirlines(airlinesclient.NewGetAirlinesParams())
			if err != nil {
				return err
			}
			if len(OK.Payload) == 0 {
				return fmt.Errorf("no airlines found")
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})
	/*	It("Should list courses", func() {
			storageURL := fmt.Sprintf("%s:%d", StorageHost, StoragePort)
			client := coursesclient.New(httptransport.New(storageURL, "", nil), strfmt.Default)
			Eventually(func() error {
				OK, err := client.GetCourses(coursesclient.NewGetCoursesParams())
				if err != nil {
					return err
				}
				if len(OK.Payload) == 0 {
					return fmt.Errorf("no courses found")
				}
				return nil
			}, "10s", "1s").ShouldNot(HaveOccurred())
		})
	*/
})
