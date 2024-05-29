package connections

type Connector interface {
	Protocol() Protocol
	User() string
	DefaultPort() int
	IsEmpty() bool
	IsValid() bool
	TestConnection(server Server) error
	Run(server Server, cmd string, exp string) error
	Open(server Server) error
	Close()
}
