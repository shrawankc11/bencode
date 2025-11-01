package bencode

import (
	"fmt"
	"reflect"
)

func encodeString(val string) []byte {
	return fmt.Appendf([]byte(""), "%d:%s", len(val), val)
}

func encodeNumber(val string) []byte {
	return fmt.Appendf([]byte(""), "i%se", val)
}

func upRes(res []byte, app []byte) []byte {
	return fmt.Appendf(res, "%s", app)
}

func Marshal(val any) ([]byte, error) {
	v := reflect.ValueOf(val)
	return marshalCore(v, reflect.TypeOf(val))
}

func marshalCore(v reflect.Value, t reflect.Type) ([]byte, error) {
	if !v.IsValid() {
		return nil, fmt.Errorf("no value passed")
	}

	var result []byte
	switch v.Kind() {
	case reflect.Slice:
		result = upRes(result, []byte("l"))
		for i := 0; i < v.Len(); i++ {
			res, _ := marshalCore(v.Index(i), v.Index(i).Type())
			result = upRes(result, res)
		}
		result = upRes(result, []byte("e"))
	case reflect.Struct:
		result = upRes(result, []byte("d"))
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			keyName := fieldType.Tag.Get("bencode")
			if len(keyName) == 0 {
				continue
			}
			res, _ := marshalCore(reflect.ValueOf(keyName), t.Field(i).Type)
			result = upRes(result, res)
			res, _ = marshalCore(field, t.Field(i).Type)
			result = upRes(result, res)
		}
		result = upRes(result, []byte("e"))
	case reflect.String:
		result = upRes(result, encodeString(v.String()))
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		result = upRes(result, encodeNumber(fmt.Sprintf("%d", v.Int())))
	case reflect.Float32, reflect.Float64:
		result = upRes(result, encodeNumber(fmt.Sprintf("%f", v.Float())))
	default:
		return nil, fmt.Errorf("type not supported %s\n", v.Kind())
	}

	return result, nil
}
