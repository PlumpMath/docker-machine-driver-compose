package brooklyn

type Application struct {
	Name     string
	Location string
	Type     string
	SshUserKey string
}


type HostAndPort struct {
	Host string
	Port int
	HasBracketlessColons bool
}

type SshHostAddress struct {
	User string
	HostAndPort HostAndPort
}


