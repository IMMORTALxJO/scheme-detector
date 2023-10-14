package schemedetector

var schemesMapping = map[string][]string{
	"ssh":        {"22"},
	"http":       {"80", "8080"},
	"https":      {"443"},
	"mysql":      {"3306"},
	"pgsql":      {"5432", "6432"},
	"memcached":  {"11211"},
	"redis":      {"6379"},
	"prometheus": {"9090"},
	"kafka":      {"9092"},
	"amqp":       {"5672"},
}

var hostHints = []string{
	"address",
	"uri",
	"url",
	"endpoint",
	"host",
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
