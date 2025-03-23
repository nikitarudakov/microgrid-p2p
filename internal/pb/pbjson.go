package pb

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func m() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		AllowPartial:  true,
		UseProtoNames: true,
	}
}

func u() protojson.UnmarshalOptions {
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
}

func ToProto[T protoreflect.ProtoMessage](in interface{}, out T) T {
	// Marshal the Model to JSON
	b, err := json.Marshal(in)
	if err != nil {
		return out
	}

	// Marshal to Proto
	if err := u().Unmarshal(b, out); err != nil {
		return out
	}

	// Return nil by default
	return out
}

func FromProto[T interface{}](in protoreflect.ProtoMessage, out T) T {
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
