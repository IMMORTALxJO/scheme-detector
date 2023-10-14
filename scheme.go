package schemedetector

import (
	"net/url"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

type Scheme struct {
	URL       *url.URL `json:"-"`
	Engine    string   `json:"engine"`
	Port      string   `json:"port"`
	Host      string   `json:"host"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	Arguments string   `json:"arguments,omitempty"`
}

func (s *Scheme) setEngine(input string) {
	s.URL.Scheme = input
	s.Engine = input
	log.Debugf("scheme: SetEngine('%s')", input)
	s.regenerate()
}

func (s *Scheme) setPort(input string) {
	if input == "" || s.Port != "" {
		return
	}
	s.Port = input
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

func (s *Scheme) setArguments(input string) {
	if input == "" || s.Arguments != "" {
		return
	}
	s.Arguments = input
	s.URL.RawQuery = s.Arguments
	log.Debugf("scheme: setArguments('%s')", input)
}

func (s *Scheme) String() string {
	return s.URL.String()
}

func (s *Scheme) setUsername(input string) {
	if input == "" || s.Username != "" {
		return
	}
	s.Username = input
	log.Debugf("scheme: SetUsername('%s')", input)
	s.regenerate()
}
func (s *Scheme) setPassword(input string) {
	if input == "" || s.Password != "" {
		log.Debug("scheme: SetPassword: empty input or password already set")
		return
	}
	s.Password = input
	log.Debugf("scheme: SetPassword('%s')", input)
	s.regenerate()
}
func (s *Scheme) setHost(input string) {
	if input == "" || s.Host != "" {
		return
	}
	s.Host = input
	log.Debugf("scheme: SetHost('%s')", input)
	s.regenerate()
}

func (s *Scheme) regenerate() {
	host := s.Host
	if host != "" && s.Port != "" {
		host = host + ":" + s.Port
	}
	s.URL.Host = host
	log.Debugf("scheme:reg host='%s'", s.URL.Host)
	if s.Username != "" {
		if s.Password == "" {
			s.URL.User = url.User(s.Username)
		} else {
			s.URL.User = url.UserPassword(s.Username, s.Password)
		}
	}
	log.Debugf("scheme:reg user='%s'", s.URL.User.String())
}

func (s *Scheme) guessMissed() {
	if s.URL.Port() == "" && s.URL.Scheme != "" {
		port := getPortFromScheme(s.URL.Scheme)
		if port != "" {
			s.Port = port
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
	return s.URL.Scheme != "" && s.Port != "" && s.Host != ""
}

func (s *Scheme) IsIP() bool {
	return s.Host != "" && govalidator.IsIP(s.Host)
}

func (s *Scheme) IsDNSName() bool {
	return s.Host != "" && govalidator.IsDNSName(s.Host)
}

func newScheme() *Scheme {
	return &Scheme{
		URL: &url.URL{},
	}
}
