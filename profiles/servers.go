package profiles

import (
	"fmt"
	"net"
	"strconv"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Server struct {
	name     string
	hostname string
	ip       net.IP
	proto    Protocol
	port     string
}

// NewServer creates a new Server with a display Name and connection Protocol to use.
func NewServer(name, proto string) (Server, error) {
	s := Server{}

	if err := s.SetName(name); err != nil {
		return s, err
	}

	if err := s.SetProtocol(proto); err != nil {
		return s, err
	}

	s.SetPort(0)
	return s, nil
}

func (s *Server) Name() string       { return s.name }
func (s *Server) Hostname() string   { return s.hostname }
func (s *Server) IP() string         { return s.ip.String() }
func (s *Server) Protocol() Protocol { return s.proto }
func (s *Server) Port() string       { return s.port }
func (s *Server) IsEmpty() bool      { return s.name == "" && s.hostname == "" }

// SetName sets the display name for Server.
func (s *Server) SetName(name string) error {
	// Add verification to prevent escape character and other exploits.
	s.name = name
	return nil
}

// SetHostname sets the hostname to use for the server. If the hostname is an IP
// it will set both the Hostname and IP to the IP.
func (s *Server) SetHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("expected valid hostname, got empty string")
	}

	if err := s.SetIP(hostname); err == nil {
		s.hostname = hostname
		return nil
	}

	if err := validate.Var(hostname, "hostname"); err != nil {
		return err
	}

	s.hostname = hostname
	return nil
}

// SetIP sets the ip address to be used for connecting to the server. If Server.hostname is set to
// an ip, this field will automatically be set. Setting this field prevents hostname lookup.
// If hostname is unset, hostname will be set to ip.
func (s *Server) SetIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("expected valid IPv4 string, got empty string")
	}

	if i := net.ParseIP(ip); i != nil {
		s.ip = i
		if s.hostname == "" {
			s.hostname = i.String()
		}
		return nil
	}

	return fmt.Errorf("ip not valid: %s", ip)
}

// SetProtocol sets the Protocol to use to connect to the server.
func (s *Server) SetProtocol(proto string) error {
	p := StringToProtocol(proto)
	if p == INVALID {
		return fmt.Errorf("invalid protocol provided: %s", proto)
	}

	s.proto = p
	return nil
}

// SetPort sets the Port to be used when connecting to the server. Empty uses Protocol default.
func (s *Server) SetPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}

	s.port = strconv.Itoa(port)
	return nil
}
