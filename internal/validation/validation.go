package validation

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	tagNameQuery = "query"
	tagNamePath  = "path"
	tagNameJSON  = "json"
)

var validate = validator.New()

func BindQuery(
	r *http.Request,
	dst any,
) error {
	queryValues := make(map[string]string)

	for key, values := range r.URL.Query() {
		if len(values) == 0 {
			continue
		}

		queryValues[key] = values[0]
	}

	if err := bindValues(
		queryValues,
		dst,
		tagNameQuery,
	); err != nil {
		return err
	}

	return Validate(dst)
}

func BindPath(
	r *http.Request,
	dst any,
) error {
	if err := bindValues(
		mux.Vars(r),
		dst,
		tagNamePath,
	); err != nil {
		return err
	}

	return Validate(dst)
}

func BindJSON(
	r *http.Request,
	dst any,
) error {
	if err := json.NewDecoder(
		r.Body,
	).Decode(dst); err != nil {
		return err
	}

	return Validate(dst)
}

func Validate(
	value any,
) error {
	if value == nil {
		return nil
	}

	return validate.Struct(value)
}

func bindValues(
	values map[string]string,
	dst any,
	tagName string,
) error {
	value := reflect.ValueOf(dst)

	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	if value.IsNil() {
		return fmt.Errorf("destination cannot be nil")
	}

	value = value.Elem()

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a struct")
	}

	return bindStruct(
		values,
		value,
		tagName,
	)
}

func bindStruct(
	values map[string]string,
	value reflect.Value,
	tagName string,
) error {
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := valueType.Field(i)

		if !field.CanSet() {
			continue
		}

		if fieldType.Anonymous {
			if err := bindStruct(
				values,
				field,
				tagName,
			); err != nil {
				return err
			}

			continue
		}

		name := getTagName(
			fieldType,
			tagName,
		)

		if name == "" {
			continue
		}

		rawValue, exists := values[name]
		if !exists {
			continue
		}

		if err := setValue(
			field,
			rawValue,
		); err != nil {
			return fmt.Errorf(
				"invalid %s parameter %q: %w",
				tagName,
				name,
				err,
			)
		}
	}

	return nil
}

func getTagName(
	field reflect.StructField,
	tagName string,
) string {
	tag := field.Tag.Get(tagName)

	if tag == "" {
		return ""
	}

	name, _, _ := strings.Cut(
		tag,
		",",
	)

	if name == "-" {
		return ""
	}

	return name
}

func setValue(
	field reflect.Value,
	value string,
) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(
				reflect.New(field.Type().Elem()),
			)
		}

		return setValue(
			field.Elem(),
			value,
		)
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		field.SetBool(parsed)

	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:

		if field.Type() == reflect.TypeOf(
			time.Duration(0),
		) {
			parsed, err := time.ParseDuration(value)
			if err != nil {
				return err
			}

			field.SetInt(int64(parsed))
			return nil
		}

		parsed, err := strconv.ParseInt(
			value,
			10,
			field.Type().Bits(),
		)
		if err != nil {
			return err
		}

		field.SetInt(parsed)

	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:

		parsed, err := strconv.ParseUint(
			value,
			10,
			field.Type().Bits(),
		)
		if err != nil {
			return err
		}

		field.SetUint(parsed)

	case reflect.Float32,
		reflect.Float64:

		parsed, err := strconv.ParseFloat(
			value,
			field.Type().Bits(),
		)
		if err != nil {
			return err
		}

		field.SetFloat(parsed)

	case reflect.Slice:
		return setSlice(
			field,
			value,
		)

	default:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			parsed, err := parseTime(value)
			if err != nil {
				return err
			}

			field.Set(reflect.ValueOf(parsed))
			return nil
		}

		return fmt.Errorf(
			"unsupported field type: %s",
			field.Type(),
		)
	}

	return nil
}

func parseTime(value string) (time.Time, error) {
	layouts := []string{
		time.DateOnly,
		time.RFC3339,
	}

	var lastErr error

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}

		lastErr = err
	}

	return time.Time{}, lastErr
}

func setSlice(
	field reflect.Value,
	value string,
) error {
	items := strings.Split(
		value,
		",",
	)

	slice := reflect.MakeSlice(
		field.Type(),
		len(items),
		len(items),
	)

	for i, item := range items {
		if err := setValue(
			slice.Index(i),
			item,
		); err != nil {
			return err
		}
	}

	field.Set(slice)

	return nil
}
