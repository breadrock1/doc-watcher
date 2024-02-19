package server

type ServerOptions struct {
	hostAddress string
	portNumber  int
}

func BuildOptions(host string, port int) *ServerOptions {
	return &ServerOptions{
		hostAddress: host,
		portNumber:  port,
	}
}
