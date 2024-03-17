package server

type Options struct {
	hostAddress string
	portNumber  int
}

func BuildOptions(host string, port int) *Options {
	return &Options{
		hostAddress: host,
		portNumber:  port,
	}
}
