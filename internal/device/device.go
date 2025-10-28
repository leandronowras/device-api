package device

import (
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type Device struct {
	id            string
	name          string
	brand         string
	state         string
	creation_time time.Time
}

const (
	StateAvailable = "available"
	StateInUse     = "in-use"
	StateInactive  = "inactive"
)

func New(name, brand string, stateOptional ...string) (*Device, error) {
	name = strings.TrimSpace(name)
	brand = strings.TrimSpace(brand)

	state := StateAvailable
	if len(stateOptional) > 0 && strings.TrimSpace(stateOptional[0]) != "" {
		state = strings.ToLower(strings.TrimSpace(stateOptional[0]))
	}

	if name == "" {
		return nil, ErrRequired("name")
	}
	if brand == "" {
		return nil, ErrRequired("brand")
	}
	if !isValidState(state) {
		return nil, ErrInvalid("state", "state must be one of: available, in-use, inactive", http.StatusBadRequest)
	}

	id, err := uuidNewString()
	if err != nil {
		return nil, &DomainError{
			Code:    "internal_error",
			Field:   "id",
			Message: "failed to generate device ID: " + err.Error(),
			HTTP:    http.StatusInternalServerError,
		}
	}

	createdAt := time.Now().UTC()

	return &Device{
		id:            id,
		name:          name,
		brand:         brand,
		state:         state,
		creation_time: createdAt,
	}, nil
}

func NewWithID(id, name, brand, state string, creationTime time.Time) (*Device, error) {
	name = strings.TrimSpace(name)
	brand = strings.TrimSpace(brand)
	state = strings.ToLower(strings.TrimSpace(state))

	if name == "" {
		return nil, ErrRequired("name")
	}
	if brand == "" {
		return nil, ErrRequired("brand")
	}
	if !isValidState(state) {
		return nil, ErrInvalid("state", "state must be one of: available, in-use, inactive", http.StatusBadRequest)
	}

	return &Device{
		id:            id,
		name:          name,
		brand:         brand,
		state:         state,
		creation_time: creationTime,
	}, nil
}

// Stub to keep this snippet standalone; swap with "github.com/google/uuid".
func uuidNewString() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (d *Device) ValidateInvariants() error {
	// Example invariant: creation_time is immutable (once set, it must not be zero)
	if d.creation_time.IsZero() {
		return ErrInvalid("creation_time", "creation_time must be set before persistence", http.StatusBadRequest)
	}
	// Example invariant: state must always be valid
	if !isValidState(d.state) {
		return ErrInvalid("state", "invalid state", http.StatusBadRequest)
	}
	return nil
}

func isValidState(s string) bool {
	switch s {
	case StateAvailable, StateInUse, StateInactive:
		return true
	default:
		return false
	}
}

// id/creation_time are server-generated and must be empty/zero when called.
func (d *Device) ValidateForCreate() error {
	if d.id != "" {
		return ErrInvalid("id", "id must be empty on create (server-generated)", http.StatusBadRequest)
	}

	if !d.creation_time.IsZero() {
		return ErrInvalid("creation_time", "creation_time must not be set on create", http.StatusBadRequest)
	}

	if strings.TrimSpace(d.name) == "" {
		return ErrRequired("name")
	}

	if strings.TrimSpace(d.brand) == "" {
		return ErrRequired("brand")
	}

	if d.state == "" {
		d.state = StateAvailable
	} else if !isValidState(d.state) {
		return ErrInvalid("state", "state must be one of: available, in-use, inactive", http.StatusBadRequest)
	}

	return nil
}

// error when conflicts with the current resource state
// for example, deleting a device that is currently "in-use".
func ErrConflict(resource, reason string) *DomainError {
	return &DomainError{
		Code:    "conflict_" + resource,
		Field:   resource,
		Message: reason,
		HTTP:    http.StatusConflict, // 409
	}
}

func (d *Device) ID() string              { return d.id }
func (d *Device) Name() string            { return d.name }
func (d *Device) Brand() string           { return d.brand }
func (d *Device) State() string           { return d.state }
func (d *Device) CreationTime() time.Time { return d.creation_time }

func (d *Device) SetName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrRequired("name")
	}
	d.name = name
	return nil
}

func (d *Device) SetBrand(brand string) error {
	brand = strings.TrimSpace(brand)
	if brand == "" {
		return ErrRequired("brand")
	}
	d.brand = brand
	return nil
}

func (d *Device) SetState(state string) error {
	state = strings.ToLower(strings.TrimSpace(state))
	if !isValidState(state) {
		return ErrInvalid("state", "state must be one of: available, in-use, inactive", http.StatusBadRequest)
	}
	d.state = state
	return nil
}
