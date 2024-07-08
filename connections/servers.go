package connections

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Server struct {
	Name     string
	Hostname string
	IP       net.IP
	Port     int
	UseIP    bool
	Connector
	Buffers
}

// NewServer creates a new Server object with a display name, port, and stores the byte.Buffers to
// be used for results and logs output. If port is set to 0, the default port for the Connector
// will always be used.
func NewServer(hostname string, port int, results, logs *bytes.Buffer) (Server, error) {
	s := Server{}
	if err := s.SetHostname(hostname); err != nil {
		return s, err
	}

	if err := s.SetPort(port); err != nil {
		return s, err
	}

	s.Buffers = NewBuffers(s.Hostname, results, logs)

	return s, nil
}

// GetIP returns Server.ip as a string.
func (s Server) GetIP() string { return s.IP.String() }

// IsEmpty returns true of Server.hostname is not set.
func (s Server) IsEmpty() bool { return s.Hostname == "" }

// IsValid retuns true if all fields needed to connect to a server are not nil or empty.
func (s Server) IsValid() bool { err := s.Validate(); return err == nil }

// IsValid retuns an error if a field needed to connect to a server is nil or empty.
func (s Server) Validate() error {
	if s.Hostname == "" {
		return errors.New("profiles.Server.Validate: hostname is empty")
	}

	if s.Results == nil {
		return errors.New("profiles.Server.Validate: Results buffer is nil")
	}

	if s.Logs == nil {
		return errors.New("profiles.Server.Validate: Logs buffer is nil")
	}

	if s.Connector == nil {
		return errors.New("profiles.Server.Validate: Connector is nil")
	}

	return s.Connector.Validate()
}

// Run passes cmd(command) and exp(expect), along with itself, on to Connector.Run to be executed.
// See Connector.Run() for more details.
func (s Server) Run(cmd, exp string) error { return s.Connector.Run(s.Buffers, cmd, exp) }

// TestConnection tries to open a connection to the server and sends an echo command to validate
// connectivity and basic access.
func (s Server) TestConnection() error {
	_, err := Pool.Open(&s)
	if err != nil {
		return fmt.Errorf("profiles.Server.TestConnection: %s", err)
	}
	return s.Connector.TestConnection(s.Buffers)
}

// GetAddr returns the host address to use, without a port. Returns "hostname" or "ip".
func (s Server) GetHostAddr() string {
	host := s.Hostname
	// Overwrite hostname with IP if Server.usIP is true. Properly set hostnames are preferred
	// but there may be times when this is not possible.
	if s.UseIP {
		host = s.GetIP()
	}

	return host
}

// GetAddr determines the address to connect to. Returns "hostname:port" or "ip:port". If port is
// set to 0, GetAddr uses protocol's default port instead.
func (s Server) GetAddr() string {
	// If Server.port is set to 0, use the Connector's default port.
	p := s.Port
	if p == 0 {
		p = s.DefaultPort()
	}

	// return "host:port"
	return net.JoinHostPort(s.GetHostAddr(), strconv.Itoa(p))
}

// SetUseIP sets the useIP field to true or false. This field is used to determine if the ip field
// should be used instead of the hostname for connecting to the server.
func (s *Server) SetUseIP(flag bool) { s.UseIP = flag }

// SetName sets the display name for the server.
func (s *Server) SetName(name string) error {
	// INCOMPLETE: Add verification to prevent escape character and other exploits.
	s.Name = name
	return nil
}

// SetHostname sets the hostname to use for the server. If the hostname is an IP it will set both
// Server.hostname and Server.ip to the IP and set Server.useIP to true.
func (s *Server) SetHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("profiles.Server.SetHostname: hostname was empty")
	}

	// If Server.SetIP is successful then hostname is an IP so set the hostname and return.
	if err := s.SetIP(hostname); err == nil {
		s.Hostname = hostname
		s.Buffers.Hostname = hostname
		return nil
	}

	// Make sure the hostname follows some type of valid format.
	if err := validate.Var(hostname, "hostname"); err != nil {
		return err
	}

	// hostname should be valid at this point so set it.
	s.Hostname = hostname
	// Make sure we change the hostname in the Buffers struct as well.
	s.Buffers.Hostname = hostname

	// If Server.name is not already set for some reason, set it to hostname.
	if s.Name == "" {
		s.Name = hostname
	}

	return nil
}

// SetIP sets the ip address to be used for connecting to the server. If Server.hostname is set to
// an ip, this field will automatically be set. Setting this field prevents hostname lookup.
// If hostname is unset, hostname will be set to ip.
func (s *Server) SetIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("profiles.Server.SetIP: expected valid IPv4 string, got empty string")
	}

	// Let the net package tell us if ip is a valid IP or not and convert it to net.IP.
	if i := net.ParseIP(ip); i != nil {
		s.IP = i
		s.UseIP = true

		// If Server.hostname is not already set for some reason, set it now.
		if s.Hostname == "" {
			// Use net.IP.String() to guarantee a properly formated IP.
			s.Hostname = i.String()
		}

		// If Server.name is not already set for some reason, set it now.
		if s.Name == "" {
			// Use net.IP.String() to guarantee a properly formated IP.
			s.Name = i.String()
		}

		return nil
	}

	return fmt.Errorf("profiles.Server.SetIP: ip not valid: %s", ip)
}

// SetPort sets the Port to be used when connecting to the server. Setting port to 0 will cause
// Connector.DefaultPort() to be used when a connection string is created.
func (s *Server) SetPort(port int) error {
	// Negative port numbers are not valid so return an error. We could use uint16 to guarantee a
	// valid port number but using int makes things easier elsewhere. Fewer conversions needed.
	if port < 0 || port > 65535 {
		return fmt.Errorf("profiles.Server.SetPort: port must be between 0 and 65535")
	}

	s.Port = port
	return nil
}

// SetConnector sets the Connector interface to be used for connecting to the server.
// MockConnector, SSHConnector, etc. server.SetConnector(&SSHConnector{})
func (s *Server) SetConnector(connector Connector) error {
	// Setting a nil Connector could get us in trouble elsewhere.
	if connector == nil {
		return errors.New("profiles.Server.SetHandler: Connector was nil")
	}

	s.Connector = connector
	return nil
}

func GetLastBufferLine(buf *bytes.Buffer) string {
	var b []string
	s := bufio.NewScanner(buf)
	for s.Scan() {
		b = append(b, s.Text())
	}

	return b[len(b)-1]
}
