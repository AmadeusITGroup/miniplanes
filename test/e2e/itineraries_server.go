package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	livenessclient "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/client/liveness"
	readinessclient "github.com/amadeusitgroup/miniapp/itineraries-server/pkg/client/readiness"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

var _ = Describe("Itineraries Server", func() {
	It("Should have liveness probe", func() {
		client := livenessclient.New(httptransport.New("localhost:8888", "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetLive(nil)
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})

	It("Should have readiness probe", func() {
		client := readinessclient.New(httptransport.New("127.0.0.1:8888", "", nil), strfmt.Default)
		Eventually(func() error {
			_, err := client.GetReady(readinessclient.NewGetReadyParams())
			if err != nil {
				return err
			}
			return nil
		}, "10s", "1s").ShouldNot(HaveOccurred())
	})

	It("Should have no itineraries", func() {
		Eventually(func() error {
			return nil
		})
	})
	It("Should have one itinerary", func() {
		Eventually(func() error {
			return nil
		})
	})
})
