package tw_trend_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTwTrend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TwTrend Suite")
}
