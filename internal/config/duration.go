package config

import (
	"encoding/json"
	"errors"
	"time"
)

// Duration is a custom type for time.Duration with JSON unmarshalling support.
type Duration time.Duration

func (d *Duration) String() string {
	return time.Duration(*d).String()
}

// UnmarshalJSON implements the json.Unmarshaler interface for Duration to correctly unmarshal time.Duration.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(value)
		return nil
	case string:
		var err error
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(duration)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
