package tw_trend_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cthulhu/tw-trend"
	"github.com/cthulhu/tw-trend/domain"
	"github.com/cthulhu/tw-trend/resource"
	"github.com/cthulhu/tw-trend/store"
)

var _ = Describe("TwTrend", func() {
	const PORT = 9000
	server := NewServer(PORT)
	go server.Run()
	BeforeEach(func() {
		resource.MAX_RESULTS = 5
		store.VolumeDir = "fixtures"
		store.FetchDaysForReports = func() []string {
			return []string{"20191010", "20191017", "20191018"}
		}
	})
	Context("Server", func() {
		Context("On ping", func() {
			It("serves helth checks", func() {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", PORT))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
			})
		})
		Context("On words", func() {
			It("serves trending", func() {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/words", PORT))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				words := domain.WordsReport{}
				decoder := json.NewDecoder(resp.Body)
				Expect(decoder.Decode(&words)).NotTo(HaveOccurred())
				Expect(words).To(BeEquivalentTo(domain.WordsReport{
					Aggregated: &domain.Aggregated{
						TokensWithCounts: []domain.TokenWithCount{
							domain.TokenWithCount{Token: "is", Count: 92},
							domain.TokenWithCount{Token: "to", Count: 82},
							domain.TokenWithCount{Token: "i", Count: 71},
							domain.TokenWithCount{Token: "van", Count: 57},
							domain.TokenWithCount{Token: "ik", Count: 54},
						},
					},
				}))
			})
		})
		Context("On words no files", func() {
			It("serves trending", func() {
				store.VolumeDir = "tmp"
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/words", PORT))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				words := domain.WordsReport{}
				decoder := json.NewDecoder(resp.Body)
				Expect(decoder.Decode(&words)).NotTo(HaveOccurred())
				Expect(words).To(BeEquivalentTo(domain.WordsReport{
					Aggregated: &domain.Aggregated{TokensWithCounts: []domain.TokenWithCount{}},
				}))
			})
		})
		Context("On hashtags", func() {
			It("serves trending", func() {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/hashtags", PORT))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				var hashtags domain.HashtagsReport
				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&hashtags)
				Expect(err).NotTo(HaveOccurred())
				Expect(hashtags).To(BeEquivalentTo(domain.HashtagsReport{
					Aggregated: &domain.Aggregated{
						TokensWithCounts: []domain.TokenWithCount{
							domain.TokenWithCount{Token: "amsterdam", Count: 15},
							domain.TokenWithCount{Token: "icdirect", Count: 12},
							domain.TokenWithCount{Token: "ade", Count: 7},
							domain.TokenWithCount{Token: "ade19", Count: 6},
							domain.TokenWithCount{Token: "amsterdamdanceevent", Count: 5},
						},
					},
				}))
			})
		})
		Context("On data", func() {
			It("serves file", func() {
				resp, err := http.Get(fmt.Sprintf("http://localhost:%d/data/20191010", PORT))
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(200))
				var tweet domain.Tweet
				decoder := json.NewDecoder(resp.Body)
				err = decoder.Decode(&tweet)
				Expect(err).NotTo(HaveOccurred())
				Expect(tweet).To(BeEquivalentTo(domain.Tweet{Text: "A story about #amsterdam", Hashtags: []string{"amsterdam"}, Tokens: []string{"about", "story"}}))
				err = decoder.Decode(&tweet)
				Expect(err).NotTo(HaveOccurred())
				Expect(tweet).To(BeEquivalentTo(domain.Tweet{Hashtags: []string{}, Tokens: []string{}}))
				err = decoder.Decode(&tweet)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func NewTestVolumeFileService() domain.VolumeFileService {
	return nil
}
