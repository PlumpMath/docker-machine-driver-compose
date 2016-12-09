package compose

// Application holds information about amp application
type Application struct {
	Name         string
	Location     string
	Type         string
	SSHUserKey   string
	OsName       string
	OsVersion    string
	Skip         bool
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

// GetSSHHostname return Host Address
func (s SSHHostAddress) GetSSHHostname() (string, error) {
	return s.HostAndPort.Host, nil
}

// GetSSHPort return SSH Port
func (s SSHHostAddress) GetSSHPort() (int, error) {
	return s.HostAndPort.Port, nil
}
