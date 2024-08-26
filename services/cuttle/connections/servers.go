package connections

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/chadeldridge/cuttle-server/db"
	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Server struct {
	ID         int64
	Name       string
	Hostname   string
	IP         net.IP
	UseIP      bool
	Connectors map[Protocol]Connector
	Buffers
}

// NewServer creates a new Server object with a display name, and stores the byte.Buffers to
// be used for results and logs output.
func NewServer(hostname string, results, logs *bytes.Buffer) (Server, error) {
	s := Server{Connectors: make(map[Protocol]Connector)}
	if err := s.SetHostname(hostname); err != nil {
		return s, err
	}

	s.Buffers = NewBuffers(s.Hostname, results, logs)

	return s, nil
}

func NewFromServerData(data db.ServerData) (Server, error) {
	s := Server{
		ID:         data.ID,
		Name:       data.Name,
		Hostname:   data.Hostname,
		IP:         net.ParseIP(data.Hostname),
		UseIP:      data.UseIP,
		Connectors: make(map[Protocol]Connector),
	}

	// If there are no connectors, return the server as is.
	if data.ConnectorIDs == "[]" {
		return s, nil
	}

	var conns []int64
	err := json.Unmarshal([]byte(data.ConnectorIDs), &conns)
	if err != nil {
		return s, fmt.Errorf("profiles.NewFromServerData: %w", err)
	}

	// INCOMPLETE: Add a way to get the connector from the database.

	return s, nil
}

func (s Server) ToServerData() db.ServerData {
	sd := db.ServerData{
		Name:     s.Name,
		Hostname: s.Hostname,
		IP:       s.IP.String(),
		UseIP:    s.UseIP,
	}

	if s.ID != 0 {
		sd.ID = s.ID
	}

	if len(s.Connectors) == 0 {
		sd.ConnectorIDs = "[]"
		return sd
	}

	conns := make([]int64, len(s.Connectors))
	for _, c := range s.Connectors {
		conns = append(conns, c.GetID())
	}

	data, _ := json.Marshal(conns)
	sd.ConnectorIDs = string(data)

	return sd
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

	if s.Connectors == nil || len(s.Connectors) == 0 {
		return nil
	}

	for p, c := range s.Connectors {
		// There should never be an invalid connector in the list.
		if p == INVALID {
			return fmt.Errorf("profiles.Server.Validate: Connector listed as invalid in Connectors")
		}

		// Make sure the connector is valid.
		err := c.Validate()
		if err != nil {
			return fmt.Errorf("profiles.Server.Validate: %w", err)
		}
	}

	return nil
}

// Open creates a connection to the server using the given Protocol. If Server.useIP is true then
// Server.ip will be used instead of Server.hostname.
func (s Server) Open(proto Protocol) error {
	if !s.UseIP && s.Hostname == "" {
		return errors.New("profiles.Server.Open: hostname is empty")
	}

	err := s.Connectors[proto].Open(s.GetHostAddr(), s.Buffers)
	if err != nil {
		return fmt.Errorf("profiles.Server.Open: %w", err)
	}

	return nil
}

// Run passes cmd(command) and exp(expect), along with itself, on to Connector.Run to be executed.
// See Connector.Run() for more details.
func (s Server) Run(proto Protocol, cmd, exp string) error {
	return s.Connectors[proto].Run(s.Buffers, cmd, exp)
}

// TestConnection tries to open a connection to the server and sends an echo command to validate
// connectivity and basic access.
func (s Server) TestConnection(proto Protocol) error {
	conn := s.Connectors[proto]
	err := conn.Open(s.GetHostAddr(), s.Buffers)
	if err != nil {
		return fmt.Errorf("profiles.Server.TestConnection: %w", err)
	}

	return conn.TestConnection(s.Buffers)
}

func (s Server) Close(proto Protocol, force bool) error {
	return s.Connectors[proto].Close(force)
}

// GetAddr returns the host address to use, without a port. Returns "hostname" or "ip" if UseIP.
func (s Server) GetHostAddr() string {
	host := s.Hostname
	// Overwrite hostname with IP if Server.usIP is true. Properly set hostnames are preferred
	// but there may be times when this is not possible.
	if s.UseIP {
		host = s.GetIP()
	}

	return host
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
		return fmt.Errorf("profiles.Server.SetHostname: hostname not valid: %s", hostname)
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

// SetConnector sets the Connector interface to be used for connecting to the server.
// MockConnector, SSHConnector, etc. server.SetConnector(&SSHConnector{})
func (s *Server) SetConnector(conn Connector) error {
	// Setting a nil Connector could get us in trouble elsewhere.
	if conn == nil {
		return errors.New("profiles.Server.SetHandler: Connector was nil")
	}

	s.Connectors[conn.Protocol()] = conn
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
