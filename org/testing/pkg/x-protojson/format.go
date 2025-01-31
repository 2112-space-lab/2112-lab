package xprotojson

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const emptyJSON = "{}"

var compact = protojson.MarshalOptions{
	Multiline:         false,
	UseProtoNames:     true,
	UseEnumNumbers:    false,
	AllowPartial:      true,
	EmitUnpopulated:   true,
	EmitDefaultValues: true,
}

var pretty = protojson.MarshalOptions{
	Indent:            "  ",
	Multiline:         true,
	UseEnumNumbers:    false,
	UseProtoNames:     true,
	AllowPartial:      true,
	EmitDefaultValues: true,
	EmitUnpopulated:   true,
}

func FormatCompact(m protoreflect.ProtoMessage) string {
	if m == nil {
		return emptyJSON
	}
	return compact.Format(m)
}

func FormatPretty(m protoreflect.ProtoMessage) string {
	if m == nil {
		return emptyJSON
	}
	return pretty.Format(m)
}
