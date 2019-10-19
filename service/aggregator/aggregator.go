package aggregator

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cthulhu/tw-trend/domain"
)

func Aggregate(reader io.Reader, label string, maxAggAmount int) (*domain.Aggregated, error) {
	decoder := json.NewDecoder(reader)
	aggregatedHash := make(map[string]int)
	var tweet domain.Tweet
	var err error
	for {
		err = decoder.Decode(&tweet)
		if err == io.EOF {
			break
		}
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
		tLow := strings.ToLower(t)
		if _, exists := aggregatedHash[tLow]; !exists {
			aggregatedHash[tLow] = 1
		} else {
			aggregatedHash[tLow]++
		}
	}
}
