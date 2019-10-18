package tokenizer

import (
	"sort"
	"strings"

	"gopkg.in/jdkato/prose.v2"
)

func Tokenize(str string) ([]string, error) {
	tokens := []string{}
	doc, err := prose.NewDocument(str)
	if err != nil {
		return tokens, err
	}
	for _, tok := range doc.Tokens() {
		if isAllowedTag(tok.Tag) && isAllowedToken(tok.Text) {
			tokens = append(tokens, strings.ToLower(tok.Text))
		}
	}
	sort.Strings(tokens)
	return tokens, err
}

var allowedTags = [...]string{"MD", "NN", "NNP", "NNPS", "NNS", "PDT", "POS", "PRP", "PRP$", "RB", "RBR", "RBS", "RP", "SYM", "TO", "UH", "VB", "VBD", "VBG", "VBN", "VBP", "VBZ", "WDT", "WP", "WP$", "WRB"}

func isAllowedTag(tag string) bool {
	for _, aTag := range allowedTags {
		if tag == aTag {
			return true
		}
	}
	return false
}

func isAllowedToken(token string) bool {
	if strings.Index(token, "@") == 0 {
		return false
	}
	if strings.Index(token, "http") == 0 {
		return false
	}
	if strings.Index(token, "#") == 0 {
		return false
	}
	return true
}
