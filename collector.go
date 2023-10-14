package schemedetector

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

func FromEnv() []*Scheme {
	envs := map[string]string{}
	for _, envStr := range os.Environ() {
		splitted := strings.SplitN(envStr, "=", 2)
		log.Debugf("env: %+v", splitted)
		envs[splitted[0]] = splitted[1]
	}
	return FromMap(envs)
}

// skipcq: GO-R1005
func FromMap(input map[string]string) []*Scheme {
	var result []*Scheme

	var procceed []string
	log.SetLevel(log.DebugLevel)

	keys := mapToKeys(filterMap(input))
	for _, k := range keys {
		log.Debugf("parse: key='%s' value='%v'", k.name, k.value)
		if stringInArray(k.name, procceed) {
			continue
		}
		item := newScheme()
		// collect all simillar keys
		group := k.findSimilars(keys)
		log.Debugf("parse: group='%v'", group)
		// search for URI in keys
		for _, k := range group {
			parsed := k.uri
			if parsed == nil {
				log.Debugf("parse: k='%v' is not URL", k)
				continue
			}
			log.Debugf("parse: k='%v' is URL %v", k, parsed.Hostname())
			item.setEngine(parsed.Scheme)
			item.setHost(parsed.Hostname())
			item.setPort(parsed.Port())
			item.setPath(parsed.Path)
			item.setUsername(parsed.User.Username())
			item.setArguments(parsed.RawQuery)
			if p, hasPass := parsed.User.Password(); hasPass {
				item.setPassword(p)
			}
			log.Debugf("procceed+='%s'", k.name)
			procceed = append(procceed, k.name)
		}
		// search parameters in other keys
		for _, k := range group {
			log.Debugf("parse:2 '%s'", k.name)
			if stringInArray(k.name, procceed) {
				log.Debugf("parse:2 '%s' already proceed", k.name)
				continue
			}
			if k.hints.host && govalidator.IsDNSName(k.value) || govalidator.IsIP(k.value) {
				item.setHost(k.value)
			} else if k.hints.port && govalidator.IsPort(k.value) {
				item.setPort(k.value)
			} else if k.hints.user {
				item.setUsername(k.value)
			} else if k.hints.pass {
				item.setPassword(k.value)
			} else if k.hints.path {
				item.setPath(k.value)
			}
		}
		item.guessMissed()
		if item.isFull() {
			for _, k := range group {
				log.Debugf("procceed+='%s'", k.name)
				procceed = append(procceed, k.name)
			}
			log.Debugf("parse: Succeed group %v", group)
			log.Debugf("parse: RESULT '%s'", item)
			result = append(result, item)
		}
	}
	return result
}

func mapToKeys(input map[string]string) []*key {
	var keys []string

	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result []*key

	for _, k := range keys {
		v := input[k]
		result = append(result, newKey(k, v))
	}
	return result
}

func filterMap(input map[string]string) map[string]string {
	excludeRegexp := os.Getenv("SCHEME_DETECTOR_EXCLUDE")
	if excludeRegexp == "" {
		return input
	}
	filter := regexp.MustCompile(excludeRegexp)
	filtered := map[string]string{}
	for k, v := range input {
		if !filter.Match([]byte(k)) {
			filtered[k] = v
		}
	}
	return filtered
}
