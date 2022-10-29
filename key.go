package schemedetector

import (
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

const splitChar = "_"

type key struct {
	name  string
	value string
}

func (k *key) findSimilars(candidates []*key) []*key {
	matches := []*key{}
	chunkedKey := strings.Split(strings.ToLower(k.name), splitChar)
	sliceLen := len(chunkedKey) - 1
	if sliceLen == 0 {
		sliceLen = 1
	}
	for _, candidate := range candidates {
		chunkedCandidate := strings.Split(strings.ToLower(candidate.name), splitChar)
		if len(chunkedCandidate) < len(chunkedKey) {
			continue
		}
		if strings.Join(chunkedCandidate[0:sliceLen], "_") == strings.Join(chunkedKey[0:sliceLen], "_") {
			matches = append(matches, candidate)
		}
	}
	log.Debugf("matches: %+v", matches)
	return matches
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
