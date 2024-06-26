package connections

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Server struct {
	name     string
	hostname string
	ip       net.IP
	port     int
	useIP    bool
	Connector
	Results *bytes.Buffer
	Logs    *bytes.Buffer
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

	s.Results = results
	s.Logs = logs

	return s, nil
}

// Name returns the value of Server.name.
func (s Server) Name() string { return s.name }

// Hostname returns the value of Server.hostname.
func (s Server) Hostname() string { return s.hostname }

// IP returns Server.ip as a string.
func (s Server) IP() string { return s.ip.String() }

// Port retuns the value of Server.port.
func (s Server) Port() int { return s.port }

// UseIP returns true if the Connector should use the set Server.ip instead of the hostname.
func (s Server) UseIP() bool { return s.useIP }

// IsEmpty returns true of Server.hostname is not set.
func (s Server) IsEmpty() bool { return s.hostname == "" }

// IsValid retuns true if all fields needed to connect to a server are not nil or empty.
func (s Server) IsValid() bool { err := s.Validate(); return err == nil }

// IsValid retuns an error if a field needed to connect to a server is nil or empty.
func (s Server) Validate() error {
	if s.hostname == "" {
		return errors.New("connections.Server.Validate: hostname is empty")
	}

	if s.Results == nil {
		return errors.New("connections.Server.Validate: Results buffer is nil")
	}

	if s.Logs == nil {
		return errors.New("connections.Server.Validate: Logs buffer is nil")
	}

	if s.Connector == nil {
		return errors.New("connections.Server.Validate: Connector is nil")
	}

	return s.Connector.Validate()
}

// Run passes cmd(command) and exp(expect), along with itself, on to Connector.Run to be executed.
// See Connector.Run() for more details.
func (s Server) Run(cmd, exp string) error { return s.Connector.Run(s, cmd, exp) }

// TestConnection tries to open a connection to the server and sends an echo command to validate
// connectivity and basic access.
func (s Server) TestConnection() error {
	_, err := Pool.Open(&s)
	if err != nil {
		return fmt.Errorf("connections.Server.TestConnection: %s", err)
	}
	return s.Connector.TestConnection(s)
}

// PrintResults adds the formated result to the Server.Results buffer.
func (s Server) PrintResults(eventTime time.Time, result string, err error) {
	if err != nil {
		fmt.Fprintf(s.Results, "%s: %s...%s: %s\n", eventTime.Format("2006/01/02 15:04:05"), s.hostname, result, err)
	}

	fmt.Fprintf(s.Results, "%s: %s...%s\n", eventTime.Format("2006/01/02 15:04:05"), s.hostname, result)
}

// Logs sends the returned connection data to the Server.Logs buffer.
func (s Server) Log(eventTime time.Time, txt string) {
	txt = strings.TrimSpace(txt)
	fmt.Fprintf(s.Logs, "%s %s@%s:~ %s\n", eventTime.Format("2006/01/02 15:04:05"), s.User(), s.Hostname(), txt)
}

// GetAddr determines the address to connect to. Returns "hostname:port" or "ip:port". If port is
// set to 0, GetAddr uses protocol's default port instead.
func (s Server) GetAddr() string {
	host := s.hostname
	// Overwrite hostname with IP if Server.usIP is true. Properly set hostnames are preferred
	// but there may be times when this is not possible.
	if s.useIP {
		host = s.IP()
	}

	// If Server.port is set to 0, use the Connector's default port.
	p := s.port
	if p == 0 {
		p = s.DefaultPort()
	}

	// return "host:port"
	return net.JoinHostPort(host, strconv.Itoa(p))
}

// SetUseIP sets the useIP field to true or false. This field is used to determine if the ip field
// should be used instead of the hostname for connecting to the server.
func (s *Server) SetUseIP(flag bool) { s.useIP = flag }

// SetName sets the display name for the server.
func (s *Server) SetName(name string) error {
	// INCOMPLETE: Add verification to prevent escape character and other exploits.
	s.name = name
	return nil
}

// SetHostname sets the hostname to use for the server. If the hostname is an IP it will set both
// Server.hostname and Server.ip to the IP and set Server.useIP to true.
func (s *Server) SetHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("connections.Server.SetHostname: hostname was empty")
	}

	// If Server.SetIP is successful then hostname is an IP so set the hostname and return.
	if err := s.SetIP(hostname); err == nil {
		s.hostname = hostname
		return nil
	}

	// Make sure the hostname follows some type of valid format.
	if err := validate.Var(hostname, "hostname"); err != nil {
		return err
	}

	// hostname should be valid at this point so set it.
	s.hostname = hostname

	// If Server.name is not already set for some reason, set it to hostname.
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
		return fmt.Errorf("connections.Server.SetIP: expected valid IPv4 string, got empty string")
	}

	// Let the net package tell us if ip is a valid IP or not and convert it to net.IP.
	if i := net.ParseIP(ip); i != nil {
		s.ip = i
		s.useIP = true

		// If Server.hostname is not already set for some reason, set it now.
		if s.hostname == "" {
			// Use net.IP.String() to guarantee a properly formated IP.
			s.hostname = i.String()
		}

		// If Server.name is not already set for some reason, set it now.
		if s.name == "" {
			// Use net.IP.String() to guarantee a properly formated IP.
			s.name = i.String()
		}

		return nil
	}

	return fmt.Errorf("connections.Server.SetIP: ip not valid: %s", ip)
}

// SetPort sets the Port to be used when connecting to the server. Setting port to 0 will cause
// Connector.DefaultPort() to be used when a connection string is created.
func (s *Server) SetPort(port int) error {
	// Negative port numbers are not valid so return an error. We could use uint16 to guarantee a
	// valid port number but using int makes things easier elsewhere. Fewer conversions needed.
	if port < 0 || port > 65535 {
		return fmt.Errorf("connections.Server.SetPort: port must be between 0 and 65535")
	}

	s.port = port
	return nil
}

// SetConnector sets the Connector interface to be used for connecting to the server.
// MockConnector, SSHConnector, etc. server.SetConnector(&SSHConnector{})
func (s *Server) SetConnector(connector Connector) error {
	// Setting a nil Connector could get us in trouble elsewhere.
	if connector == nil {
		return errors.New("connections.Server.SetHandler: Connector was nil")
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
