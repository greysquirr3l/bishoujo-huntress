// Package huntress provides a client for the Huntress API
package huntress

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// Helper functions for service implementations

// addQueryParams adds the parameters in params as URL query parameters to s.
func addQueryParams(path string, params interface{}) (string, error) {
	v, err := encodeURLValues(params)
	if err != nil {
		return path, err
	}

	if len(v) > 0 {
		return fmt.Sprintf("%s?%s", path, v), nil
	}

	return path, nil
}

// extractPagination extracts pagination information from HTTP response headers
func extractPagination(resp *http.Response) *Pagination {
	if resp == nil {
		return nil
	}

	pagination := &Pagination{}

	if currentPage := resp.Header.Get("X-Page"); currentPage != "" {
		if val, err := parseInt(currentPage); err == nil {
			pagination.CurrentPage = val
		}
	}

	if perPage := resp.Header.Get("X-Per-Page"); perPage != "" {
		if val, err := parseInt(perPage); err == nil {
			pagination.PerPage = val
		}
	}

	if totalPages := resp.Header.Get("X-Total-Pages"); totalPages != "" {
		if val, err := parseInt(totalPages); err == nil {
			pagination.TotalPages = val
		}
	}

	if totalItems := resp.Header.Get("X-Total-Items"); totalItems != "" {
		if val, err := parseInt(totalItems); err == nil {
			pagination.TotalItems = val
		}
	}

	return pagination
}

// parseInt parses a string to an integer, returning 0 if parsing fails
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// encodeURLValues converts a struct into URL-encoded query parameters.
// It uses the `url` tag from field declarations.
func encodeURLValues(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "", nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("query: can only encode structs, got %T", v)
	}

	values := url.Values{}
	err := addValues(values, val)
	if err != nil {
		return "", err
	}

	return values.Encode(), nil
}

// addValues adds the values from the struct to the specified url.Values.
func addValues(values url.Values, val reflect.Value) error {
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Handle embedded structs
		if field.Anonymous {
			if fieldValue.Kind() == reflect.Struct {
				err := addValues(values, fieldValue)
				if err != nil {
					return err
				}
			}
			continue
		}

		tag := field.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}

		name, opts := parseTag(tag)
		if name == "" {
			continue
		}

		if opts.Contains("omitempty") && isEmptyValue(fieldValue) {
			continue
		}

		// Handle different types
		var strValues []string

		switch fieldValue.Kind() {
		case reflect.String:
			strValues = []string{fieldValue.String()}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValues = []string{fmt.Sprintf("%d", fieldValue.Int())}
		case reflect.Bool:
			if fieldValue.Bool() {
				strValues = []string{"true"}
			} else {
				strValues = []string{"false"}
			}
		case reflect.Slice:
			if fieldValue.Len() == 0 {
				continue
			}

			if fieldValue.Type().Elem().Kind() == reflect.String {
				for j := 0; j < fieldValue.Len(); j++ {
					strValues = append(strValues, fieldValue.Index(j).String())
				}
			} else {
				// For other slice types, convert to JSON
				b, err := json.Marshal(fieldValue.Interface())
				if err != nil {
					return err
				}
				strValues = []string{string(b)}
			}
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				b, err := json.Marshal(fieldValue.Interface())
				if err != nil {
					return err
				}
				strValues = []string{string(b)}
			}
		default:
			// For other types, convert to JSON
			b, err := json.Marshal(fieldValue.Interface())
			if err != nil {
				return err
			}
			strValues = []string{string(b)}
		}

		for _, v := range strValues {
			values.Add(name, v)
		}
	}
	return nil
}

// parseTag splits a struct field's url tag into its name and options.
func parseTag(tag string) (string, tagOptions) {
	idx := 0
	for idx < len(tag) {
		if tag[idx] == ',' {
			return tag[:idx], tagOptions(tag[idx+1:])
		}
		idx++
	}
	return tag, tagOptions("")
}

// tagOptions represents the options specified in a struct field's url tag.
type tagOptions string

// Contains checks whether the options include the specified opt.
func (o tagOptions) Contains(opt string) bool {
	if len(o) == 0 {
		return false
	}

	s := string(o)
	for s != "" {
		var next string
		i := 0
		for i < len(s) && s[i] != ',' {
			i++
		}
		if i < len(s) {
			next = s[i+1:]
			s = s[:i]
		}
		if s == opt {
			return true
		}
		s = next
	}
	return false
}

// isEmptyValue checks if a value is empty according to Go's notion of empty.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
