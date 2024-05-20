package connections

import (
	"fmt"
)

type MockHandler struct {
	proto Protocol
	user  string
}

var MockDefaultPort = 555

func (h MockHandler) User() string     { return h.user }
func (h MockHandler) DefaultPort() int { return MockDefaultPort }
func (h MockHandler) IsEmpty() bool    { return h.proto == INVALID && h.user == "" }
func (h MockHandler) IsValid() bool    { return h.proto == MOCK && h.user != "" }
func (h MockHandler) TestConnection(server Server) error {
	res, err := h.Run(server, "echo", "any")
	if err != nil {
		return err
	}

	fmt.Fprintf(server.Results, "%s...%s", server.Hostname(), res)
	return nil
}

func (h MockHandler) Run(server Server, cmd, expect string) (string, error) {
	h.Log(server, "mock ok")
	return "ok", nil
}

func (h MockHandler) Log(server Server, txt string) {
	fmt.Fprintf(server.Logs, "%s@%s:~ %s", h.user, server.Hostname(), txt)
}
