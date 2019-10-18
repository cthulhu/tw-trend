package aggregator

import (
	"fmt"
	"sort"

	"github.com/cthulhu/tw-trend/domain"
)

func Aggregate(tweets []domain.Tweet, label string, maxAggAmount int) (*domain.Aggregated, error) {
	aggregatedHash := make(map[string]int)
	for _, tweet := range tweets {
		switch label {
		case "words":
			increment(aggregatedHash, tweet.Tokens)
		case "hashtags":
			increment(aggregatedHash, tweet.Hashtags)
		default:
			return nil, fmt.Errorf("unknown label")
		}
	}

	aggregated := &domain.Aggregated{TokensWithCounts: []domain.TokenWithCount{}}

	for token, count := range aggregatedHash {
		aggregated.TokensWithCounts = append(aggregated.TokensWithCounts, domain.TokenWithCount{token, count})
	}
	sort.Slice(aggregated.TokensWithCounts, func(i, j int) bool {
		return aggregated.TokensWithCounts[i].Count > aggregated.TokensWithCounts[j].Count
	})
	if len(aggregated.TokensWithCounts) > maxAggAmount {
		aggregated.TokensWithCounts = aggregated.TokensWithCounts[0:maxAggAmount]
	}
	return aggregated, nil
}

func increment(aggregatedHash map[string]int, tokens []string) {
	for _, t := range tokens {
		if _, exists := aggregatedHash[t]; !exists {
			aggregatedHash[t] = 1
		} else {
			aggregatedHash[t]++
		}
	}
}
