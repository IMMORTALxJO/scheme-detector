package main

import (
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

var hostHints = []string{
	"host",
	"hostname",
	"address",
	"endpoint",
}
var pathHints = []string{
	"path",
	"name",
	"db",
}
var passHints = []string{
	"pass",
	"password",
	"pwn",
}
var userHints = []string{
	"user",
	"username",
}
var portHints = []string{
	"port",
}

var hints = []string{
	"host",
	"hostname",
	"address",
	"endpoint",
	"path",
	"name",
	"db",
	"pass",
	"password",
	"pwn",
	"user",
	"username",
	"port",
}

func FromEnv() []*Scheme {
	envs := map[string]string{}
	for _, envStr := range os.Environ() {
		splitted := strings.SplitN(envStr, "=", 1)
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
		item := NewScheme()
		// collect all simillar keys
		group := k.findSimilars(keys, hints)
		log.Debugf("parse: group='%v'", group)
		log.Debugf("parse: item='%v'", item)
		// search for URI in keys
		for _, k := range group {
			parsed := k.getURL()
			if parsed == nil {
				log.Debugf("parse: k='%v' is not URL", k)
				continue
			}
			log.Debugf("parse: k='%v' is URL %v", k, parsed.Hostname())
			item.SetEngine(parsed.Scheme)
			item.SetHost(parsed.Hostname())
			item.SetPort(parsed.Port())
			item.SetPath(parsed.Path)
			item.SetUsername(parsed.User.Username())
			if p, hasPass := parsed.User.Password(); hasPass {
				item.SetPassword(p)
			}
			procceed = append(procceed, k.name)
		}
		// search parameters in other keys
		for _, k := range group {
			log.Debugf("parse:2 '%s'", k.name)
			if stringInArray(k.name, procceed) {
				continue
			}
			if k.hasHints(hostHints) && (govalidator.IsDNSName(k.value) || govalidator.IsIP(k.value)) {
				item.SetHost(k.value)
			}
			if k.hasHints(portHints) && govalidator.IsNumeric(k.value) {
				item.SetPort(k.value)
			}
			if k.hasHints(userHints) {
				item.SetUsername(k.value)
			}
			if k.hasHints(passHints) {
				item.SetPassword(k.value)
			}
			if k.hasHints(pathHints) {
				item.SetPath(k.value)
			}
		}
		item.Guess()
		if item.IsFull() {
			for _, k := range group {
				procceed = append(procceed, k.name)
			}
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
