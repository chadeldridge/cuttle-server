package connections

type Handler interface {
	TestConnection() error
	Run(cmd string, expect string) (string, error)
}
