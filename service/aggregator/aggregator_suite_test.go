package aggregator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAggregator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aggregator Suite")
}
