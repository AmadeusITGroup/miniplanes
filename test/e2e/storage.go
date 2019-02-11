package e2e

import (
	"fmt"

	airlinesclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/airlines"
	airportsclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/airports"
	coursesclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/courses"
	livenessclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/liveness"
	readinessclient "github.com/amadeusitgroup/miniapp/storage/pkg/gen/client/readiness"
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
	It("Should list courses", func() {
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
})
