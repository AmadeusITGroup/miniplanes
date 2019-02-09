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
	goflag "flag"
	"os"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/spf13/pflag"
	//"k8s.io/client-go/rest"
)

func TestE2E(t *testing.T) {
	RunE2ETests(t)
}

var (
	//kubeConfig         *rest.Config
	kubeConfigFilePath string
	StorageHost        string
	StoragePort        int
)

func TestMain(m *testing.M) {

	//pflag.StringVar(&kubeConfigFilePath, "kubeconfig", "", "Path to kubeconfig")
	pflag.StringVar(&StorageHost, "storage-host", "", "storage host/service name")
	pflag.IntVar(&StoragePort, "storage-port", 0, "storage port number")
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()
	goflag.CommandLine.Parse([]string{})

	os.Exit(m.Run())
}

// RunE2ETests runs e2e test
func RunE2ETests(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "miniapp suite")
}
