package domainenum

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

// PotentialHandlerState is exposed as an intermediate type for validating if input is a valid enum value or not
type PotentialHandlerState string

// handlerStateStr internal type holding the actual string representation
type handlerStateStr string

// HandlerState is safe wrapper around valid enum values
type HandlerState struct{ source handlerStateStr }

// Validate checks the potential enum value
func (m PotentialHandlerState) Validate() (HandlerState, error) {
	return HandlerStates.FromString(m)
}

// String return underlying enum value or "UNKNOWN"
func (m HandlerState) String() string {
	if m.source == "" {
		return "UNKNOWN"
	}
	return string(m.source)
}

func (m HandlerState) upperString() string { return strings.ToUpper(string(m.String())) }

type handlerStates struct{}

// HandlerStates is a reference object holding helpers for allowed values
var HandlerStates = handlerStates{}

// Unknown is helper for default fallback value
func (dm handlerStates) Unknown() HandlerState { return HandlerState{""} }

// Started helper for value
func (dm handlerStates) Started() HandlerState { return HandlerState{"Started"} }

// Completed helper for value
func (dm handlerStates) Completed() HandlerState { return HandlerState{"Completed"} }

// Failed helper for value
func (dm handlerStates) Failed() HandlerState { return HandlerState{"Failed"} }

// FromString checks potential enum value
func (dm handlerStates) FromString(s PotentialHandlerState) (HandlerState, error) {
	switch strings.ToUpper(string(s)) {
	case dm.Started().upperString():
		return dm.Started(), nil
	case dm.Completed().upperString():
		return dm.Completed(), nil
	case dm.Failed().upperString():
		return dm.Failed(), nil
	}
	return dm.Unknown(), errors.New("unknown handler state [" + string(s) + "]")
}

// MarshalJSON override for enum value to json
func (m HandlerState) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// UnmarshalJSON override for json value to enum
func (m *HandlerState) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		m.source = HandlerStates.Unknown().source
		return nil
	}

	var value PotentialHandlerState
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	s, err := value.Validate()
	if err != nil {
		return err
	}
	m.source = s.source
	return nil
}
