package schemedetector

import (
	"os"
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

func FromMap(input map[string]string) []*Scheme {
	result := []*Scheme{}
	procceed := []string{}
	keys := mapToKeys(input)
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
			parsed := k.getURL()
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
			procceed = append(procceed, k.name)
		}
		// search parameters in other keys
		for _, k := range group {
			log.Debugf("parse:2 '%s'", k.name)
			if stringInArray(k.name, procceed) {
				log.Debugf("parse:2 '%s' already proceed", k.name)
				continue
			}
			if k.hasHints(hostHints) && govalidator.IsDNSName(k.value) || govalidator.IsIP(k.value) {
				item.setHost(k.value)
			} else if k.hasHints(portHints) && govalidator.IsNumeric(k.value) {
				item.setPort(k.value)
			} else if k.hasHints(userHints) {
				item.setUsername(k.value)
			} else if k.hasHints(passHints) {
				item.setPassword(k.value)
			} else if k.hasHints(pathHints) {
				item.setPath(k.value)
			}
		}
		item.guessMissed()
		if item.isFull() {
			for _, k := range group {
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
	result := []*key{}
	for k, v := range input {
		result = append(result, &key{k, v})
	}
	return result
}
