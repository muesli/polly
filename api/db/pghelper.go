package db

import (
	"database/sql/driver"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
)

// StringSlice is a custom string slice for postgres encoding
type StringSlice []string

// BigintSlice is a custom int64 slice for postgres encoding
type BigintSlice []int64

// Scan converts a postgres value into a string slice
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	asBytes, ok := value.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	str := string(asBytes)

	// convert escapes (\" -> "" for CSV parser) and remove braces
	r := strings.NewReplacer(`\"`, `""`, `\\`, `\`)
	str = r.Replace(str[1 : len(str)-1])

	csvReader := csv.NewReader(strings.NewReader(str))

	slice, err := csvReader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	(*s) = StringSlice(slice)

	return nil
}

// Value constructs a single postgres query string from a slice of strings
func (s StringSlice) Value() (driver.Value, error) {
	// string escapes.
	// \ => \\
	// " => \"
	r := strings.NewReplacer(`"`, `\"`, `\`, `\\`)

	t := []string{}
	for _, elem := range s {
		t = append(t, `"`+r.Replace(elem)+`"`)
	}
	return "{" + strings.Join(t, ",") + "}", nil
}

// Scan converts a postgres value into an int64 slice
func (s *BigintSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	asBytes, ok := value.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	str := string(asBytes)

	if len(str) <= 2 {
		return nil
	}

	// Remove braces
	str = str[1 : len(str)-1]

	slice := strings.Split(str, ",")

	for _, elem := range slice {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			return err
		}
		(*s) = append(*s, i)
	}

	return nil
}

// Value constructs a single postgres query string from a slice of int64
func (s BigintSlice) Value() (driver.Value, error) {
	t := []string{}

	for _, elem := range s {
		t = append(t, strconv.FormatInt(elem, 10))
	}
	return "{" + strings.Join(t, ",") + "}", nil
}
