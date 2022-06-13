package main

import (
	"net/url"
	"strings"
)

type key struct {
	name  string
	value string
}

type hintedKey struct {
	name     string
	original *key
	matched  bool
	replacer string
	matches  []string
}

/*
Replace strings from 'hints' on 'replace_by'.
Example:
		KEYS_HINT_LIST = {'port','engine', 'database'}
foo("BACKEND_PORT") => { "result": "BACKEND_XXX", "matches": ["port"] }
foo("DATABASE_ENGINE") => { "result": "XXX_XXX", "matches": ["database", "host"] }
foo("BACKEND_PORT") => { "result": "BACKEND_XXX", "matches": ["port"] }
foo("NGINX_HOST") => { "result": "NGINX_HOST", "matches": [] }
*/
func (k *key) replaceHints(hints []string) *hintedKey {
	splitChar := "_"
	replaceBy := "XXX"
	result := []string{}
	matches := []string{}
	for _, chunk := range strings.Split(strings.ToLower(k.name), splitChar) {
		if stringInArray(chunk, hints) {
			result = append(result, replaceBy)
			if !stringInArray(chunk, matches) {
				matches = append(matches, chunk)
			}
		} else {
			result = append(result, chunk)
		}
	}
	return &hintedKey{
		name:     strings.Join(result, splitChar),
		original: k,
		matches:  matches,
		replacer: replaceBy,
		matched:  len(matches) > 0,
	}
}

/*
Check if string has any special predefined word in splitted view
KEYS_HINT_LIST = {'port','engine', 'database'}
foo("BACKEND_PORT") => True
foo("DATABASE_ENGINE") => True
foo("BACKEND_PORT") => True
foo("NGINX_HOST") => False
*/
func (k *key) hasHints(hints []string) bool {
	return k.replaceHints(hints).matched
}

/*
Search for similar strings based on similarity of splitted 'chunks'
Example:
foo(
		"DATABASE_HOST",
		["NGINX_HOST", "DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASS"]
) => ["DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASS"]
*/

func (k *key) findSimilars(candidates []*key, hints []string) []*key {
	replaced := k.replaceHints(hints)
	var result []*key
	for _, candidate := range candidates {
		if replaced.name == candidate.replaceHints(hints).name {
			result = append(result, candidate)
		}
	}
	return result
}

func (k *key) getURL() *url.URL {
	parsed, err := url.Parse(k.value)
	if err != nil || parsed.Host == "" {
		return nil
	}
	return parsed
}

func (k *key) String() string {
	return k.name
}
