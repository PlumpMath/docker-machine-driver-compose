package compose

import (
	"errors"
	"strconv"
	"strings"
)

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
// type HostAndPort struct {
// 	 Host                 string
//	 Port                 int
//	 HasBracketlessColons bool
// }

// SSHHostAddress SSH Host Address
type SSHHostAddress struct {
	User        string
	HostAndPort string
}

// GetSSHHostname return Host Address
func (s SSHHostAddress) GetSSHHostname() (string, error) {
	tokens := strings.Split(s.HostAndPort, ":")
	return tokens[0], nil
}

// GetSSHPort return SSH Port
func (s SSHHostAddress) GetSSHPort() (int, error) {
	tokens := strings.Split(s.HostAndPort, ":")
	if len(tokens) != 2 {
		return 22, errors.New("Invalid HostAndPort value")
	}

	return strconv.Atoi(tokens[1])
}
