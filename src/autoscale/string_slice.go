package autoscale

import (
	"database/sql/driver"
	"strings"
)

// StringSlice is a slice of strings
type StringSlice []string

// Value converts a string slice to something the driver value can handle. In this case,
// it creates a CSV.
func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

// Scan converts a DB value back into a StringSlice.
func (s *StringSlice) Scan(src interface{}) error {
	u8 := src.([]uint8)
	ba := make([]byte, 0, len(u8))
	for _, b := range u8 {
		ba = append(ba, byte(b))
	}

	str := string(ba)
	*s = strings.Split(str, ",")
	return nil
}
