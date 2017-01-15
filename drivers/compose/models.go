package compose

// Application holds information about amp application
type Application struct {
	Name         string
	Location     string
	Type         string
	SSHUserKey   string
	OsName       string   `json:",omitempty"`
	OsVersion    string   `json:",omitempty"`
	Skip         bool     `json:",omitempty"`
	TemplateSize string   `json:",omitempty"`
	OpenPorts    []string `json:",omitempty"`
	NewRelic     bool     `json:",omitempty"`
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
