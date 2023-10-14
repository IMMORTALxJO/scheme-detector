package schemedetector

import (
	"net/url"
	"strings"
)

const splitChar = "_"

type key struct {
	name    string
	value   string
	chunked []string
	masked  string
	uri     *url.URL
	hints   *keyHints
}

type keyHints struct {
	host bool
	path bool
	port bool
	user bool
	pass bool
}

func (k *key) findSimilars(candidates []*key) []*key {
	var result []*key
	for _, candidate := range candidates {
		if candidate.masked == k.masked {
			result = append(result, candidate)
		}
	}
	return result
}

func newKey(k string, v string) *key {
	// check if value is URL
	uri, err := url.Parse(v)
	if err != nil || uri.Host == "" {
		uri = nil
	}
	// split key by splitChar
	chunked := strings.Split(strings.ToLower(k), splitChar)

	// mask all hints
	hostHinted := false
	portHinted := false
	userHinted := false
	passHinted := false
	pathHinted := false
	var masked []string
	for _, chunk := range chunked {
		mask := false
		if stringInArray(chunk, hostHints) {
			mask = true
			hostHinted = true
		} else if stringInArray(chunk, portHints) {
			mask = true
			portHinted = true
		} else if stringInArray(chunk, userHints) {
			mask = true
			userHinted = true
		} else if stringInArray(chunk, passHints) {
			mask = true
			passHinted = true
		} else if stringInArray(chunk, pathHints) {
			mask = true
			pathHinted = true
		}
		if mask {
			masked = append(masked, "XXX")
		} else {
			masked = append(masked, chunk)
		}
	}

	return &key{
		name:    k,
		value:   v,
		chunked: chunked,
		masked:  strings.Join(masked, splitChar),
		uri:     uri,
		hints: &keyHints{
			host: hostHinted,
			path: pathHinted,
			port: portHinted,
			user: userHinted,
			pass: passHinted,
		},
	}
}

func (k *key) String() string {
	return k.name
}
