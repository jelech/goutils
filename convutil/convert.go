// Package convert provides type conversion utilities.
package convutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ToString converts various types to string
func ToString(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts various types to int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ToInt64 converts various types to int64
func ToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ToFloat64 converts various types to float64
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToBool converts various types to bool
func ToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case nil:
		return false, nil
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float32:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// ToTime converts various types to time.Time
func ToTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case nil:
		return time.Time{}, nil
	case time.Time:
		return v, nil
	case string:
		// Try common time formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time string: %s", v)
	case int64:
		// Assume Unix timestamp
		return time.Unix(v, 0), nil
	default:
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", value)
	}
}

// ToStringSlice converts various types to []string
func ToStringSlice(value interface{}) ([]string, error) {
	if value == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("value is not a slice or array")
	}

	result := make([]string, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = ToString(rv.Index(i).Interface())
	}

	return result, nil
}

// ToIntSlice converts various types to []int
func ToIntSlice(value interface{}) ([]int, error) {
	if value == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("value is not a slice or array")
	}

	result := make([]int, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val, err := ToInt(rv.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("cannot convert element at index %d: %w", i, err)
		}
		result[i] = val
	}

	return result, nil
}

// ToJSON converts a value to JSON string
func ToJSON(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON converts a JSON string to a value
func FromJSON(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}

// ToMap converts a struct to map[string]interface{}
func ToMap(value interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("value is not a struct")
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Get field name (use json tag if available)
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if idx := len(jsonTag); idx > 0 && jsonTag[0:1] != "," {
				if commaIdx := len(jsonTag); commaIdx > 0 {
					for j, char := range jsonTag {
						if char == ',' {
							commaIdx = j
							break
						}
					}
					fieldName = jsonTag[:commaIdx]
				} else {
					fieldName = jsonTag
				}
			}
		}

		result[fieldName] = fieldValue.Interface()
	}

	return result, nil
}

// Must functions that panic on error

// MustToString converts to string or panics
func MustToString(value interface{}) string {
	return ToString(value)
}

// MustToInt converts to int or panics
func MustToInt(value interface{}) int {
	result, err := ToInt(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToInt64 converts to int64 or panics
func MustToInt64(value interface{}) int64 {
	result, err := ToInt64(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToFloat64 converts to float64 or panics
func MustToFloat64(value interface{}) float64 {
	result, err := ToFloat64(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToBool converts to bool or panics
func MustToBool(value interface{}) bool {
	result, err := ToBool(value)
	if err != nil {
		panic(err)
	}
	return result
}

// MustToTime converts to time.Time or panics
func MustToTime(value interface{}) time.Time {
	result, err := ToTime(value)
	if err != nil {
		panic(err)
	}
	return result
}
