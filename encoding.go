package juice_cache

import (
	"encoding"
	"encoding/json"
)

// jsonSerializer implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler.
// It marshals and unmarshal the value to and from json.
type jsonSerializer struct {
	v any
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (jbm *jsonSerializer) MarshalBinary() ([]byte, error) {
	return json.Marshal(jbm.v)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (jbm *jsonSerializer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, jbm.v)
}

// jsonMarshalBinaryWrap return an encoding.BinaryMarshaler that marshals the value to json.
func jsonMarshalBinaryWrap(v any) encoding.BinaryMarshaler {
	return &jsonSerializer{v: v}
}

// jsonUnmarshalBinaryWrap return an encoding.BinaryUnmarshaler that unmarshals the value from json.
func jsonUnmarshalBinaryWrap(v any) encoding.BinaryUnmarshaler {
	return &jsonSerializer{v: v}
}
