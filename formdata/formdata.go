// Package formdata unmarshals HTML form data [url.Values] into structs.
//
// # Supported data types
//
//   - bool: Valid true values are "on", "true", "yes", "checked" and "1".
//   - int variants.
//   - uint variants.
//   - float variants.
//   - string.
//   - []byte slice: Decodes base64 strings.
//   - [encoding.TextUnmarshaler].
//   - Slices of the above types.
//   - Pointers to the above types.
//
// # Struct tags
//
// The tag name is "formdata". The first argument is the name of the form key.
// It defaults to the name of the field if not given.
// If the name is "-", the field is skipped.
// The second argument specifies if the field is required.
// Required fields will give an error if they are missing in the form data.
//
// # Data validation
//
// Package formdata does not validate values other than converting types.
// Types can validate themselves by implementing the [encoding.TextUnmarshaler] interface.
//
// Any errors that occur during unmarshaling are accumulated and returned as an instance of [Errors].
//
// # Default values
//
// You can set default values by setting the value of non-required fields before calling [Unmarshal].
//
// # Example
//
//	type Form struct {
//		Name     string `formdata:"name,required"`
//		FavColor string `formdata:"favcolor"`
//		Age      int    `formdata:"age"`
//		Address  string `formdata:"-"`
//	}
//	values := url.Values{}
//	values.Set("name", "Gopher")
//	values.Set("age", "42")
//	var form Form
//	form.FavColor = "blue"
//	formdata.Unmarshal(values, &form)
package formdata

import (
	"encoding"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	byteSliceType           = reflect.TypeOf((*[]byte)(nil)).Elem()
	textUnmarshalerType     = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	errNotAPointerToStruct  = errors.New("formdata: expected pointer to struct")
	errRequiredFieldMissing = errors.New("required but missing")
)

var rxbool = regexp.MustCompile(`(?i)^(1|checked|on|true|yes)$`)

func isTextUnmarshaler(vtype reflect.Type) bool {
	return reflect.PointerTo(vtype).Implements(textUnmarshalerType)
}

func decodeScalar(v reflect.Value, vtype reflect.Type, data string) error {
	// special case for TextUnmarshaler
	if v.CanAddr() && isTextUnmarshaler(vtype) {
		return v.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(data))
	}
	switch vtype.Kind() {
	case reflect.Bool:
		v.SetBool(rxbool.MatchString(data))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(data, 10, vtype.Bits())
		if err != nil {
			return err
		}
		v.SetInt(x)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(data, 10, vtype.Bits())
		if err != nil {
			return err
		}
		v.SetUint(x)
		return nil
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(data, vtype.Bits())
		if err != nil {
			return err
		}
		v.SetFloat(x)
		return nil
	case reflect.String:
		v.SetString(data)
		return nil
	}
	// special case for []byte slice
	if vtype.ConvertibleTo(byteSliceType) {
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(decoded))
		return nil
	}
	return fmt.Errorf("cannot unmarshal into unsupported type %q", vtype)
}

func decodeSlice(vtype reflect.Type, vals []string) (reflect.Value, error) {
	if elem := vtype.Elem(); !isTextUnmarshaler(elem) {
		// fast path for []string
		if elem.Kind() == reflect.String {
			return reflect.ValueOf(slices.Clone(vals)), nil
		}
		// special case for []byte slice
		if len(vals) == 1 && vtype.ConvertibleTo(byteSliceType) {
			decoded, err := base64.StdEncoding.DecodeString(vals[0])
			if err != nil {
				return reflect.Zero(vtype), err
			}
			return reflect.ValueOf(decoded), nil
		}
	}
	vnew := reflect.New(vtype).Elem()
	vtype = vtype.Elem()
	tmpval := reflect.New(vtype)
	for _, data := range vals {
		if err := decodeScalar(tmpval.Elem(), vtype, data); err != nil {
			return vnew, err
		}
		vnew = reflect.Append(vnew, tmpval.Elem())
	}
	return vnew, nil
}

func decodeValue(v reflect.Value, vtype reflect.Type, vals []string) error {
	// deref pointer
	for vtype.Kind() == reflect.Pointer {
		vtype = vtype.Elem()
		v.Set(reflect.New(vtype))
		v = v.Elem()
	}
	// decode slices
	if vtype.Kind() == reflect.Slice {
		vnew, err := decodeSlice(vtype, vals)
		if err != nil {
			return err
		}
		v.Set(vnew)
		return nil
	}
	// decode scalar
	return decodeScalar(v, vtype, vals[0])
}

func decodeField(data url.Values, v reflect.Value, field reflect.StructField) (Error, bool) {
	if field.Anonymous || !field.IsExported() {
		return Error{}, true
	}
	tag := field.Tag.Get("formdata")
	name, rest, _ := strings.Cut(tag, ",")
	if name == "-" { // skip
		return Error{}, true
	}
	if name == "" { // default name
		name = field.Name
	}
	vals, exists := data[name]
	if !exists { // missing value
		if rest == "required" {
			return Error{
				Field: name,
				Err:   errRequiredFieldMissing,
			}, false
		}
		return Error{}, true
	}
	fv := v.FieldByIndex(field.Index)
	if len(vals) == 0 { // empty value
		fv.Set(reflect.Zero(field.Type))
		return Error{}, true
	}
	if err := decodeValue(fv, field.Type, vals); err != nil {
		return Error{
			Field: name,
			Value: vals,
			Err:   err,
		}, false
	}
	return Error{}, true
}

func structTypeOf(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() != reflect.Pointer {
		return nil
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil
	}
	return t
}

// Unmarshal decodes a map into a struct.
// Parameter v must be a pointer to a struct.
// Returns [Errors] if there are any unmarshaling errors.
func Unmarshal(data map[string][]string, v any) error {
	structType := structTypeOf(reflect.TypeOf(v))
	if structType == nil {
		return errNotAPointerToStruct
	}
	var errors Errors
	structValue := reflect.ValueOf(v).Elem()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if err, ok := decodeField(data, structValue, field); !ok {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Error holds the information of an unmarshaling error.
type Error struct {
	// The field where the error occurred.
	Field string
	// The input value that was passed.
	Value []string
	// The actual error.
	Err error
}

// Error implements the error interface.
func (err Error) Error() string {
	return fmt.Sprintf("error in field %q: %s", err.Field, err.Err)
}

// Errors holds multiple unmarshaling errors.
type Errors []Error

// Error implements the error interface.
func (errs Errors) Error() string {
	switch len(errs) {
	case 0:
		return ""
	case 1:
		return errs[0].Error()
	default:
		var b []byte
		b = append(b, errs[0].Error()...)
		for _, err := range errs[1:] {
			b = append(b, '\n')
			b = append(b, err.Error()...)
		}
		return string(b)
	}
}
