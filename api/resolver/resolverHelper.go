package resolver

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// M global marshaler
func m() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		AllowPartial:  true,
		UseProtoNames: true,
	}
}

// FromProto takes a proto object and converts to a swagger model
func fromProto[T interface{}](in protoreflect.ProtoMessage, out T) T {
	// Marshal to JSON
	b, err := m().Marshal(in)
	if err != nil {
		return out
	}

	// Unmarshal into the Out Object
	if err := json.Unmarshal(b, &out); err != nil {
		return out
	}

	// Return nil by default
	return out
}
