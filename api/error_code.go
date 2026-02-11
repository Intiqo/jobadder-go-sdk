package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// UnmarshalJSON implements a custom JSON unmarshaler for ErrorCode that handles
// both string and numeric values. The JobAdder API may return numeric error codes
// in some responses, while the OpenAPI spec defines them as string enums.
func (e *ErrorCode) UnmarshalJSON(data []byte) error {
	// Try string first (expected case from the spec).
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*e = ErrorCode(s)
		return nil
	}

	// Fall back to number (observed in some API error responses).
	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		*e = ErrorCode(strconv.FormatFloat(n, 'f', -1, 64))
		return nil
	}

	return fmt.Errorf("ErrorCode: cannot unmarshal %s", string(data))
}
