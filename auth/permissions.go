package auth

import (
	"encoding/json"
	"fmt"
)

type Permissions struct {
	perms map[string]bool
}

var (
	validMethods = []string{"POST", "GET", "PUT", "DELETE"}
	defaultPerms = map[string]bool{
		"POST":   false, // Create
		"GET":    false, // Read
		"PUT":    false, // Update
		"DELETE": false, // Delete
	}
)

var ErrInvalidMethod = fmt.Errorf("invalid method")

func NewPermissions() Permissions {
	return Permissions{perms: defaultPerms}
}

func ValidMethods() []string {
	return validMethods
}

func IsValidMethod(method string) bool {
	for _, k := range validMethods {
		if k == method {
			return true
		}
	}

	return false
}

func (p Permissions) Set(method string, value bool) error {
	if !IsValidMethod(method) {
		return fmt.Errorf("auth.Permissions.Set: %w", ErrInvalidMethod)
	}

	p.perms[method] = value
	return nil
}

func (p Permissions) HasPermission(method string) bool {
	if !IsValidMethod(method) {
		return false
	}

	return p.perms[method]
}

func (p Permissions) Get(method string) (bool, error) {
	if !IsValidMethod(method) {
		return false, fmt.Errorf("auth.Permissions.Get: %w", ErrInvalidMethod)
	}

	return p.perms[method], nil
}

func (p Permissions) AllowPost()   { p.perms["POST"] = true }
func (p Permissions) AllowGet()    { p.perms["GET"] = true }
func (p Permissions) AllowPut()    { p.perms["PUT"] = true }
func (p Permissions) AllowDelete() { p.perms["DELETE"] = true }

func (p Permissions) CanCreate() bool { return p.perms["POST"] }
func (p Permissions) CanRead() bool   { return p.perms["GET"] }
func (p Permissions) CanUpdate() bool { return p.perms["PUT"] }
func (p Permissions) CanDelete() bool { return p.perms["DELETE"] }

func (p Permissions) DenyPost()   { p.perms["POST"] = false }
func (p Permissions) DenyGet()    { p.perms["GET"] = false }
func (p Permissions) DenyPut()    { p.perms["PUT"] = false }
func (p Permissions) DenyDelete() { p.perms["DELETE"] = false }

func (p Permissions) AllowAll() {
	for k := range p.perms {
		p.perms[k] = true
	}
}

func (p Permissions) DenyAll() {
	for k := range p.perms {
		p.perms[k] = false
	}
}

func (p Permissions) Marshal() ([]byte, error) {
	// We only need to write values that are true. This will save on space.
	t := make(map[string]bool)
	for k, v := range p.perms {
		if v {
			t[k] = v
		}
	}

	data, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("auth.Permissions.Marshal: %w", err)
	}

	return data, nil
}

func (p Permissions) Unmarshal(data []byte) error {
	t := make(map[string]bool)
	err := json.Unmarshal(data, &t)
	if err != nil {
		return fmt.Errorf("auth.Permissions.Unmarshal: %w", err)
	}

	for k, v := range t {
		p.perms[k] = v
	}

	return nil
}

func UnmarshalPermissions(data []byte) (Permissions, error) {
	p := NewPermissions()
	err := p.Unmarshal(data)
	if err != nil {
		return Permissions{}, fmt.Errorf("auth.UnmarshalPermissions: %w", err)
	}

	return p, nil
}
