package adapter_test

import (
	"fmt"
	"log"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"

	"gopkg.in/yaml.v2"

	"github.com/dherbric/test-service-adapter/adapter"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("manifest-generator", func() {

	falsy := false
	params := serviceadapter.GenerateManifestParams{

		Plan: serviceadapter.Plan{
			InstanceGroups: []serviceadapter.InstanceGroup{
				{Name: "nginx",
					AZs:                []string{"az1"},
					Instances:          1,
					Networks:           []string{"default"},
					PersistentDiskType: "default",
					VMType:             "default",
				},
			},
			Update: &serviceadapter.Update{
				Canaries:        1,
				MaxInFlight:     1,
				Serial:          &falsy,
				CanaryWatchTime: "1000-60000",
				UpdateWatchTime: "1000-60000",
			},
		},
		ServiceDeployment: serviceadapter.ServiceDeployment{
			DeploymentName: "nginx-abc8237834673",
			Releases: serviceadapter.ServiceReleases{
				{
					Name:    "nginx",
					Version: "latest",
					Jobs:    []string{"nginx"},
				},
			},
			Stemcells: []serviceadapter.Stemcell{{
				OS:      "ubuntu-trusty",
				Version: "latest",
			}},
		},
	}

	logger := log.New(GinkgoWriter, "test-service-adapter-test-logger", log.LstdFlags)
	mg := adapter.TestServiceManifestGenerator{
		Logger: logger,
	}

	Context(".GenerateManifest", func() {
		It("returns an empty manifest", func() {

			output, err := mg.GenerateManifest(params)
			yml, _ := yaml.Marshal(output)
			fmt.Println(string(yml))

			// Empty Response
			Expect(err).NotTo(HaveOccurred())
			Expect(output).ToNot(BeNil())

		})
	})

	Context(".GenerateManifest", func() {
		It("returns a valid manifest", func() {

			output, err := mg.GenerateManifest(params)
			yml, _ := yaml.Marshal(output)
			fmt.Println(string(yml))

			// Empty Response
			Expect(err).NotTo(HaveOccurred())
			Expect(output).ToNot(BeNil())

			Expect(output.Manifest.InstanceGroups).To(HaveLen(1))

			nginxInstanceGroup := output.Manifest.InstanceGroups[0]
			Expect(nginxInstanceGroup.Name).To(BeIdenticalTo("nginx"))

			Expect(nginxInstanceGroup.AZs).NotTo(HaveLen(0))

			Expect(nginxInstanceGroup.Jobs).To(HaveLen(1))

			nginxJob := nginxInstanceGroup.Jobs[0]
			Expect(nginxJob.Release).To(BeIdenticalTo("nginx"))

			Expect(nginxJob.Properties).To(HaveLen(2))

			Expect(output.Manifest.Name).To(BeIdenticalTo("nginx-abc8237834673"))
			Expect(output.Manifest.Releases).To(HaveLen(1))

			Expect(output.Manifest.Stemcells).To(HaveLen(1))

			Expect(*output.Manifest.Update).To(BeEquivalentTo(bosh.Update{
				Canaries:        1,
				MaxInFlight:     1,
				Serial:          &falsy,
				CanaryWatchTime: "1000-60000",
				UpdateWatchTime: "1000-60000",
			}))
		})
	})

})
