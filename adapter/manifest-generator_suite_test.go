package adapter_test


import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServiceBroker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Broker Suite")
}