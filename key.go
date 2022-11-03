package schemedetector

import (
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

const splitChar = "_"

type key struct {
	name    string
	value   string
	chunked []string
}

func (k *key) findSimilars(candidates []*key) []*key {
	firstMatch := []*key{}
	for _, candidate := range candidates {
		if k.isSimilar(candidate) {
			firstMatch = append(firstMatch, candidate)
		}
	}
	log.Debugf("firstMatch: %+v", firstMatch)
	weights := []int{}
	for _, candidate := range firstMatch {
		weight := 0
		for _, j := range firstMatch {
			if candidate.isSimilar(j) {
				weight = weight + 1
			}
		}
		weights = append(weights, weight)
	}
	log.Debugf("weights: %+v", weights)
	secondMatch := []*key{}
	if len(firstMatch) > 2 {
		for i, weight := range weights {
			if weight > 2 {
				secondMatch = append(secondMatch, firstMatch[i])
			}
		}
	} else {
		secondMatch = firstMatch
	}
	log.Debugf("secondMatch: %+v", secondMatch)
	return secondMatch
}

func (k key) isSimilar(candidate *key) bool {
	if k.name == candidate.name {
		return true
	}
	chunked := k.getChunked()
	candidateChunked := candidate.getChunked()
	for _, chunk := range chunked {
		if stringInArray(chunk, candidateChunked) {
			return true
		}
	}
	return false
}

func (k *key) getChunked() []string {
	k.chunked = strings.Split(strings.ToLower(k.name), splitChar)
	return k.chunked
}

func (k *key) String() string {
	return k.name
}

func (k *key) getURL() *url.URL {
	parsed, err := url.Parse(k.value)
	if err != nil || parsed.Host == "" {
		return nil
	}
	return parsed
}

func (k *key) hasHints(hints []string) bool {
	chunkedKey := strings.Split(strings.ToLower(k.name), splitChar)
	for _, hint := range hints {
		if stringInArray(hint, chunkedKey) {
			return true
		}
	}
	return false
}

func newKey(k string, v string) *key {
	return &key{
		k,
		v,
		[]string{},
	}
}
