package tokenizer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cthulhu/tw-trend/service/tokenizer"
)

var _ = Describe("Tokenizer", func() {
	var tokens []string
	var err error

	Context("On tokenize", func() {
		It("returns a set of tokens on valid strings", func() {
			tokens, err = tokenizer.Tokenize("String of words")
			Expect(tokens).To(BeEquivalentTo([]string{"string", "words"}))
			Expect(err).ShouldNot(HaveOccurred())

			tokens, err = tokenizer.Tokenize("With femkeroobol bartolomeo_ltd and Some more at the #neworder #concert in #Amsterdam. Brilliant! @ AFAS Live https://t.co/XRmwacLFHd")
			Expect(tokens).To(BeEquivalentTo([]string{"afas", "bartolomeo_ltd", "brilliant", "femkeroobol", "live"}))
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("returns nothing on not valid strings", func() {
			tokens, err = tokenizer.Tokenize("1 2 3")
			Expect(tokens).To(BeEquivalentTo([]string{}))
			Expect(err).ShouldNot(HaveOccurred())

			tokens, err = tokenizer.Tokenize("#Amsterdam @ https://t.co/XRmwacLFHd")
			Expect(tokens).To(BeEquivalentTo([]string{}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
