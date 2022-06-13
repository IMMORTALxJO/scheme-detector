package main

var schemesMapping = map[string][]string{
	"ssh":        []string{"22"},
	"http":       []string{"80", "8080"},
	"https":      []string{"443"},
	"mysql":      []string{"3306"},
	"postgres":   []string{"5432", "6432"},
	"memcached":  []string{"11211"},
	"redis":      []string{"6379"},
	"prometheus": []string{"9090"},
	"kafka":      []string{"9092"},
	"amqp":       []string{"5672"},
}

func getPortFromScheme(scheme string) string {
	if v, ok := schemesMapping[scheme]; ok {
		return v[0]
	}
	return ""
}

func getSchemeFromPort(port string) string {
	for scheme, ports := range schemesMapping {
		if stringInArray(port, ports) {
			return scheme
		}
	}
	return ""
}
