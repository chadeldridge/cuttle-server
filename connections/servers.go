package connections

import (
	"bytes"
	"errors"
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
	port     string
	useIP    bool
	Connector
	Results *bytes.Buffer
	Logs    *bytes.Buffer
}

// NewServer creates a new Server with a display Name and connection Protocol to use. If port is
// set to 0, the default port for the protocol will be used.
func NewServer(hostname string, port int, results, logs *bytes.Buffer) (Server, error) {
	s := Server{}
	if err := s.SetHostname(hostname); err != nil {
		return s, err
	}

	if err := s.SetPort(port); err != nil {
		return s, err
	}

	s.Results = results
	s.Logs = logs

	return s, nil
}

func NewPlacholderServer() Server { return Server{name: "Empty Server Profile"} }

func (s Server) Name() string     { return s.name }
func (s Server) Hostname() string { return s.hostname }
func (s Server) IP() string       { return s.ip.String() }
func (s Server) Port() string     { return s.port }
func (s Server) UseIP() bool      { return s.useIP }
func (s Server) IsEmpty() bool    { return s.hostname == "" }

func (s Server) TestConnection() error        { return s.Connector.TestConnection(s) }
func (s Server) Run(cmd, expect string) error { return s.Connector.Run(s, cmd, expect) }

func (s Server) PrintResults(result string, err error) {
	if err != nil {
		fmt.Fprintf(s.Results, "%s...%s: %s", s.hostname, result, err)
	}

	fmt.Fprintf(s.Results, "%s...%s", s.hostname, result)
}

// GetAddr determines the address to connect to. Returns "hostname:port" or "ip:port". If port is
// set to 0, GetAddr uses protocol's default port instead.
func (s Server) GetAddr() string {
	host := s.hostname
	if s.useIP {
		host = s.IP()
	}

	return net.JoinHostPort(host, s.port)
}

// SetUseIP sets the useIP field to true or false. This field is used to determine if the ip field
// should be used instead of the hostname.
func (s *Server) SetUseIP(flag bool) { s.useIP = flag }

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
	if s.name == "" {
		s.name = hostname
	}

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
		s.useIP = true
		if s.hostname == "" {
			s.hostname = i.String()
		}
		return nil
	}

	return fmt.Errorf("ip not valid: %s", ip)
}

// SetPort sets the Port to be used when connecting to the server. Empty uses Protocol default.
func (s *Server) SetPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("port must be between 0 and 65535")
	}

	s.port = strconv.Itoa(port)
	return nil
}

// SetHandler sets the Handler interface to be used for connecting to the server.
func (s *Server) SetHandler(handler Connector) error {
	if handler == nil {
		return errors.New("could not set handler, provided handler was nil")
	}
	if handler.IsEmpty() {
		return errors.New("could not set handler, provided handler was empty")
	}

	s.Connector = handler

	if s.port == "0" {
		s.port = strconv.Itoa(handler.DefaultPort())
	}

	return nil
}
