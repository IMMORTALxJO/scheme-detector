package main

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Scheme struct {
	url      *url.URL
	port     string
	host     string
	username string
	password string
}

func (s *Scheme) SetEngine(input string) {
	if input == "" || s.url.Scheme != "" {
		return
	}
	s.url.Scheme = input
	log.Debugf("scheme: SetEngine('%s')", input)
	s.regenerate()
}

func (s *Scheme) SetPort(input string) {
	if input == "" || s.port != "" {
		return
	}
	s.port = input
	log.Debugf("scheme: SetPort('%s')", input)
	s.regenerate()
}

func (s *Scheme) SetPath(input string) {
	if input == "" || s.url.Path != "" {
		return
	}
	s.url.Path = input
	log.Debugf("scheme: SetPath('%s')", input)
}

func (s *Scheme) String() string {
	return s.url.String()
}

func (s *Scheme) SetUsername(input string) {
	if input == "" || s.username != "" {
		return
	}
	s.username = input
	log.Debugf("scheme: SetUsername('%s')", input)
	s.regenerate()
}
func (s *Scheme) SetPassword(input string) {
	if input == "" || s.password != "" {
		return
	}
	s.password = input
	log.Debugf("scheme: SetPassword('%s')", input)
	s.regenerate()
}
func (s *Scheme) SetHost(input string) {
	if input == "" || s.host != "" {
		return
	}
	s.host = input
	log.Debugf("scheme: SetHost('%s')", input)
	s.regenerate()
}

func (s *Scheme) regenerate() {
	host := s.host
	if host != "" && s.port != "" {
		host = host + ":" + s.port
	}
	s.url.Host = host
	log.Debugf("scheme:reg host='%s'", s.url.Host)
	if s.username != "" {
		if s.password == "" {
			s.url.User = url.User(s.username)
		} else {
			s.url.User = url.UserPassword(s.username, s.password)
		}
	}
	log.Debugf("scheme:reg user='%s'", s.url.User.String())
}

func (s *Scheme) Guess() {
	if s.url.Port() == "" && s.url.Scheme != "" {
		port := getPortFromScheme(s.url.Scheme)
		if port != "" {
			s.port = port
		}
	}
	if s.url.Port() != "" && s.url.Scheme == "" {
		scheme := getSchemeFromPort(s.url.Port())
		if scheme != "" {
			s.url.Scheme = scheme
		}
	}
	s.regenerate()
}

func (s *Scheme) IsFull() bool {
	return s.url.Scheme != "" && s.port != "" && s.host != ""
}

func NewScheme() *Scheme {
	return &Scheme{
		url: &url.URL{},
	}
}
