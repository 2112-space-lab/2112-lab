package fx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValueOption(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{ // cases
		{
			name:  "int value option should hasValue true with same value as input",
			input: 42,
		},
		{
			name:  "str value option should hasValue true with same value as input",
			input: "hello",
		},
		{
			name: "struct value option should hasValue true with same value as input",
			input: struct {
				v1 string
				v2 int
			}{
				v1: "hello",
				v2: 42,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := NewValueOption(tt.input)
			assert.Equal(t, tt.input, opt.Value)
			assert.True(t, opt.HasValue)
		})
	}
}

func TestNewEmptyOption(t *testing.T) {
	emptyInt := NewEmptyOption[int]()
	assert.False(t, emptyInt.HasValue)
	emptyStr := NewEmptyOption[string]()
	assert.False(t, emptyStr.HasValue)
	emptyStruct := NewEmptyOption[struct {
		v1 string
		v2 int
	}]()
	assert.False(t, emptyStruct.HasValue)
}

type Abc struct {
	T1 Option[string]
	T2 Option[string]
}

func TestOptionFormat(t *testing.T) {
	xx := NewEmptyOption[string]()
	aa := NewValueOption(Abc{
		T1: xx,
		T2: NewValueOption("abc"),
	})
	// s2 := fmt.Sprintf("%+v", aa.T1)
	s := fmt.Sprintf("%#v", aa)
	s2 := fmt.Sprintf("%v", aa)
	assert.NotEmpty(t, s2)
	assert.NotEmpty(t, s)
}
