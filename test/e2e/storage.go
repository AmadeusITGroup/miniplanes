/*
Copyright 2018 Amadeus SaS All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package e2e

import (
	"fmt"

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
