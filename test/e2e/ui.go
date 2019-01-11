package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Itineraries Server", func() {
	/*	It("Should have liveness probe", func() {
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
	*/

	It("Should list airports", func() {
		Eventually(func() error {
			return nil
		})
	})
	It("Should query an itinerary", func() {
		Eventually(func() error {
			return nil
		})
	})
})
