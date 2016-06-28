package compose

type Application struct {
	Name         string
	Location     string
	Type         string
	SshUserKey   string
	OsName       string
	OsVersion    string
	TemplateSize string
	OpenPorts    []string
}

func NewApplication() *Application {
	return &Application{}
}

type HostAndPort struct {
	Host                 string
	Port                 int
	HasBracketlessColons bool
}

type SshHostAddress struct {
	User        string
	HostAndPort HostAndPort
}
