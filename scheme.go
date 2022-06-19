package schemedetector

import (
	"net/url"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

type Scheme struct {
	URL      *url.URL
	port     string
	host     string
	username string
	password string
}

func (s *Scheme) setEngine(input string) {
	if input == "" || s.URL.Scheme != "" {
		return
	}
	s.URL.Scheme = input
	log.Debugf("scheme: SetEngine('%s')", input)
	s.regenerate()
}

func (s *Scheme) setPort(input string) {
	if input == "" || s.port != "" {
		return
	}
	s.port = input
	log.Debugf("scheme: SetPort('%s')", input)
	s.regenerate()
}

func (s *Scheme) setPath(input string) {
	if input == "" || s.URL.Path != "" {
		return
	}
	s.URL.Path = input
	log.Debugf("scheme: SetPath('%s')", input)
}

func (s *Scheme) String() string {
	return s.URL.String()
}

func (s *Scheme) setUsername(input string) {
	if input == "" || s.username != "" {
		return
	}
	s.username = input
	log.Debugf("scheme: SetUsername('%s')", input)
	s.regenerate()
}
func (s *Scheme) setPassword(input string) {
	if input == "" || s.password != "" {
		return
	}
	s.password = input
	log.Debugf("scheme: SetPassword('%s')", input)
	s.regenerate()
}
func (s *Scheme) setHost(input string) {
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
	s.URL.Host = host
	log.Debugf("scheme:reg host='%s'", s.URL.Host)
	if s.username != "" {
		if s.password == "" {
			s.URL.User = url.User(s.username)
		} else {
			s.URL.User = url.UserPassword(s.username, s.password)
		}
	}
	log.Debugf("scheme:reg user='%s'", s.URL.User.String())
}

func (s *Scheme) guess() {
	if s.URL.Port() == "" && s.URL.Scheme != "" {
		port := getPortFromScheme(s.URL.Scheme)
		if port != "" {
			s.port = port
		}
	}
	if s.URL.Port() != "" && s.URL.Scheme == "" {
		scheme := getSchemeFromPort(s.URL.Port())
		if scheme != "" {
			s.URL.Scheme = scheme
		}
	}
	s.regenerate()
}

func (s *Scheme) isFull() bool {
	return s.URL.Scheme != "" && s.port != "" && s.host != ""
}

func (s *Scheme) IsIP() bool {
	return s.host != "" && govalidator.IsIP(s.host)
}

func (s *Scheme) IsDNSName() bool {
	return s.host != "" && govalidator.IsDNSName(s.host)
}

func newScheme() *Scheme {
	return &Scheme{
		URL: &url.URL{},
	}
}
