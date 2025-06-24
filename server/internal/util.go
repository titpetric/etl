package internal

import (
	"encoding/json"
)

// Marshal is a "MustMarshal" version where a panic is issued if
// an invalid type to marshal has been marshalled. It should be
// used for known types, to declutter error handling.
func Marshal(in any) []byte {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return b
}

// Marshal is MustMarshalIndent with preset to as to not require
// additional arguments for that. Suitable for console output.
func MarshalIndent(in any) []byte {
	b, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		panic(err)
	}
	return b
}
