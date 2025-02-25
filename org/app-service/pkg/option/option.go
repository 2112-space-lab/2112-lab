package fx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetOrDefault generic functional type to explicitly handle both existing values or not
func GetOrDefault[T any](opt Option[T], defaultValue T) T {
	if opt.HasValue {
		return opt.Value
	}
	return defaultValue
}

// Option generic functional type to explicitly handle both existing values or not
type Option[T any] struct {
	Value    T
	HasValue bool
}

// NewValueOption initializes an Option with an actual inner value
func NewValueOption[T any](value T) Option[T] {
	return Option[T]{
		Value:    value,
		HasValue: true,
	}
}

// NewEmptyOption initializes an option without an inner value
func NewEmptyOption[T any]() Option[T] {
	return Option[T]{}
}

// OptionApply apply given function to option value if any and return option with new value. Otherwise return an empty option
func OptionApply[T any, U any](inputOption Option[T], applyFn func(inValue T) (U, error)) (Option[U], error) {
	if !inputOption.HasValue {
		return NewEmptyOption[U](), nil
	}
	newValue, err := applyFn(inputOption.Value)
	if err != nil {
		return NewValueOption(newValue), err
	}
	return NewValueOption(newValue), nil
}

// Format implements Formatter interface
func (o Option[T]) Format(w fmt.State, v rune) {
	hasHash := w.Flag('#')
	if hasHash {
		_, _ = fmt.Fprint(w, o.GoString())
	} else {
		_, _ = fmt.Fprint(w, o.String())
	}
}

// String implements Stringer interface
func (o Option[T]) String() string {
	if o.HasValue {
		return fmt.Sprintf("Some(%+v)", o.Value)
	}
	return "None"
}

// GoString implements GoStringer interface
func (o Option[T]) GoString() string {
	if o.HasValue {
		return fmt.Sprintf("Some[%T](%#v)", o.Value, o.Value)
	}
	return fmt.Sprintf("None[%T]", o.Value)
}

var _ fmt.Formatter = Option[string]{}

// MarshalJSON to marshal the json data globally
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.HasValue {
		return json.Marshal(o.Value)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON to unmarshal globally
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		o.HasValue = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = value
	o.HasValue = true
	return nil
}

// AsOption generic functional type to explicitly handle both existing values or not
func AsOption[T any](value *T) Option[T] {
	if value != nil {
		return NewValueOption(*value)
	}
	return NewEmptyOption[T]()
}

// ConvertOption converts an fx.Option[T] to fx.Option[U] using a provided transformation function.
func ConvertOption[T any, U any](input Option[T], transform func(T) U) Option[U] {
	if !input.HasValue {
		return NewEmptyOption[U]()
	}
	return NewValueOption(transform(input.Value))
}

// ConvertPointer converts a pointer *T to a new object U using a provided mapper function.
// If the pointer is nil, it returns an empty Option[U].
func ConvertPointer[T any, U any](ptr *T, mapper func(T) U) Option[U] {
	if ptr == nil {
		return NewEmptyOption[U]()
	}
	return NewValueOption(mapper(*ptr))
}

// ConvertToFloatOption *float64 to fx.Option[float64]
func ConvertToFloatOption(f *float64) Option[float64] {
	if f == nil {
		return NewEmptyOption[float64]()
	}
	return NewValueOption(*f)
}

// ConvertToFloatPtr converts fx.Option[float64] to *float64
func ConvertToFloatPtr(opt Option[float64]) *float64 {
	if !opt.HasValue {
		return nil
	}
	return &opt.Value
}

// ConvertToStrPtr converts fx.Option[string] to *string
func ConvertToStrPtr(opt Option[string]) *string {
	if !opt.HasValue {
		return nil
	}
	return &opt.Value
}
