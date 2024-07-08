package connections

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type Buffers struct {
	User     string
	Hostname string
	Results  *bytes.Buffer
	Logs     *bytes.Buffer
}

// NewBuffers creates a new Buffers object with the hostname, Results buffer, and Logs buffer set.
// You will need to set the User field before using the Buffers object.
func NewBuffers(hostname string, results, logs *bytes.Buffer) Buffers {
	return Buffers{
		Hostname: hostname,
		Results:  results,
		Logs:     logs,
	}
}

// Clear resets both the Results and Logs buffers.
func (b Buffers) Clear() { b.Results.Reset(); b.Logs.Reset() }

// ClearResults resets the Results buffer.
func (b Buffers) ClearResults() { b.Results.Reset() }

// ClearLogs resets the Logs buffer.
func (b Buffers) ClearLogs() { b.Logs.Reset() }

// PrintResults adds the formated result to the Server.Results buffer.
func (b Buffers) PrintResults(eventTime time.Time, result string, err error) {
	if err != nil {
		fmt.Fprintf(b.Results, "%s: %s...%s: %s\n", eventTime.Format("2006/01/02 15:04:05"), b.Hostname, result, err)
		return
	}

	fmt.Fprintf(b.Results, "%s: %s...%s\n", eventTime.Format("2006/01/02 15:04:05"), b.Hostname, result)
}

// Logs sends the returned connection data to the Server.Logs buffer.
func (b Buffers) Log(eventTime time.Time, txt string) {
	txt = strings.TrimSpace(txt)
	fmt.Fprintf(b.Logs, "%s %s@%s:~ %s\n", eventTime.Format("2006/01/02 15:04:05"), b.User, b.Hostname, txt)
}
