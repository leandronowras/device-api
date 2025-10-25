package device

import (
	"testing"
	"time"
)

func TestDeviceCreationAndImmutability(t *testing.T) {
	type want struct {
		state       string
		setRecently bool
	}

	type in struct {
		name  string
		brand string
		state string // optional; empty => expect default
	}

	cases := []struct {
		name     string
		input    in
		want     want
		wantErr  bool
		errCode  string
		errField string
	}{
		{
			name:  "defaults to available when state omitted",
			input: in{name: "iPhone 15", brand: "Apple", state: ""},
			want:  want{state: StateAvailable, setRecently: true},
		},
		{
			name:  "accepts provided valid state (inactive)",
			input: in{name: "ThinkPad X1", brand: "Lenovo", state: StateInactive},
			want:  want{state: StateInactive, setRecently: true},
		},
		{
			name:     "rejects empty name",
			input:    in{name: "   ", brand: "Apple"},
			wantErr:  true,
			errCode:  "required",
			errField: "name",
		},
		{
			name:     "rejects empty brand",
			input:    in{name: "PS5", brand: "   "},
			wantErr:  true,
			errCode:  "required",
			errField: "brand",
		},
		{
			name:     "rejects invalid state",
			input:    in{name: "Router", brand: "TP-Link", state: "broken"},
			wantErr:  true,
			errCode:  "invalid_state",
			errField: "state",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var (
				dev *Device
				err error
			)

			if tt.input.state == "" {
				dev, err = New(tt.input.name, tt.input.brand)
			} else {
				dev, err = New(tt.input.name, tt.input.brand, tt.input.state)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if de, ok := err.(*DomainError); ok {
					if de.Code != tt.errCode {
						t.Fatalf("want error code %q, got %q (msg=%q)", tt.errCode, de.Code, de.Message)
					}
					if de.Field != tt.errField {
						t.Fatalf("want error field %q, got %q", tt.errField, de.Field)
					}
					t.Logf("pretty message to user: %s", de.Message)
				} else {
					t.Fatalf("expected *DomainError, got %T: %v", err, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Success assertions
			if dev == nil {
				t.Fatalf("device is nil")
			}
			if dev.id == "" {
				t.Fatalf("id should be generated")
			}
			if dev.state != tt.want.state {
				t.Fatalf("want state %q, got %q", tt.want.state, dev.state)
			}
			if tt.want.setRecently {
				// creation_time should be within a small recent window
				if time.Since(dev.creation_time) > 2*time.Second || time.Since(dev.creation_time) < 0 {
					t.Fatalf("creation_time not set recently: now=%s creation_time=%s",
						time.Now().UTC().Format(time.RFC3339Nano),
						dev.creation_time.Format(time.RFC3339Nano))
				}
			}

			t.Logf("created device id=%s state=%s at=%s", dev.id, dev.state, dev.creation_time.Format(time.RFC3339))
		})
	}
}
