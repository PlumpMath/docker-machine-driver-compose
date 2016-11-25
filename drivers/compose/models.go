package compose

// Application holds information about amp application
type Application struct {
	Name         string
	Location     string
	Type         string
	SSHUserKey   string
	TemplateSize string
	OpenPorts    []string
}

// NewApplication return empty Application
func NewApplication() *Application {
	return &Application{}
}

// HostAndPort information
type HostAndPort struct {
	Host                 string
	Port                 int
	HasBracketlessColons bool
}

// SSHHostAddress SSH Host Address
type SSHHostAddress struct {
	User        string
	HostAndPort HostAndPort
}
