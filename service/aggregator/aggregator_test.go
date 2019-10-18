package aggregator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cthulhu/tw-trend/domain"
	. "github.com/cthulhu/tw-trend/service/aggregator"
)

var _ = Describe("Aggregator", func() {
	var aggregated *domain.Aggregated
	var err error
	maxAggregatedAmount := 2
	Context("On aggregate", func() {
		It("counts words", func() {
			tweets := []domain.Tweet{
				domain.Tweet{Tokens: []string{"one", "two", "one"}},
				domain.Tweet{Tokens: []string{"three", "four", "two", "one"}},
			}
			aggregated, err = Aggregate(tweets, "words", maxAggregatedAmount)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(aggregated).To(BeEquivalentTo(&domain.Aggregated{
				[]domain.TokenWithCount{
					domain.TokenWithCount{Token: "one", Count: 3},
					domain.TokenWithCount{Token: "two", Count: 2},
				},
			}))
		})
		It("counts hashtags", func() {
			tweets := []domain.Tweet{
				domain.Tweet{Hashtags: []string{"one", "two", "one"}},
				domain.Tweet{Hashtags: []string{"three", "two", "three"}},
				domain.Tweet{Hashtags: []string{"four", "two", "one"}},
			}
			aggregated, err = Aggregate(tweets, "hashtags", maxAggregatedAmount)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(aggregated).To(BeEquivalentTo(&domain.Aggregated{
				[]domain.TokenWithCount{
					domain.TokenWithCount{Token: "one", Count: 3},
					domain.TokenWithCount{Token: "two", Count: 3},
				},
			}))
		})
	})
})
