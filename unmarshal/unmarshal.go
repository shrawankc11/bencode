package unmarshal

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func convertSlice(val reflect.Value, dest any) error {
	destVal := reflect.ValueOf(dest)
	newSlice := reflect.MakeSlice(destVal.Elem().Type(), 0, val.Len()+1)
	for i := 0; i < val.Len(); i++ {
		res := val.Index(i)
		if res.CanConvert(destVal.Elem().Type().Elem()) {
			res = res.Convert(destVal.Elem().Type().Elem())
			newSlice = reflect.Append(newSlice, res)
		}
	}
	destVal.Elem().Set(newSlice)
	return nil
}

type CoreRes struct {
	RefVal reflect.Value
	Val    any
}

func UnMarshal(e []byte, val any) (err error, v any) {
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return fmt.Errorf("function expects a apointer received value"), nil
	}

	read := 0
	valRef := reflect.ValueOf(val)
	corRes := CoreRes{
		Val: val,
		RefVal: valRef,
	}
	err, res := unMarshalCore(e, corRes, &read)

	if err != nil {
		return err, nil
	}

	if valRef.Elem().Kind() != res.RefVal.Kind() {
		return fmt.Errorf("type mismatched"), nil
	}

	// convertSlice(resRef, val)
	valRef.Elem().Set(res.RefVal)

	return nil, res.Val
}

func unMarshalCore(e []byte, val CoreRes, i *int) (error, *CoreRes) {
	var err error

	reader := bytes.NewReader(e[*i:])
	initialByte := make([]byte, 1)
	_, err = reader.Read(initialByte)

	if err != nil {
		return err, nil
	}

	ibStr := string(initialByte)

	switch {
	case ibStr >= "0" && ibStr <= "9":
		str := ""
		for ; e[*i] != ':'; *i++ {
			str += string(e[*i])
		}
		skip, err := strconv.Atoi(str)
		if err != nil {
			log.Fatal(err)
			return err, nil
		}

		// incr for ":"
		*i++
		strVal := string(e[*i : *i+skip])
		// NOTE
		*i += skip

		return nil, &CoreRes{RefVal: reflect.ValueOf(strVal), Val: strVal}

	case ibStr == "i":
		*i++
		strData := ""
		for ; e[*i] != 'e'; *i++ {
			strData += string(e[*i])
		}
		*i++

		if i := strings.IndexByte(strData, '.'); i != -1 {
			f, err := strconv.ParseFloat(strData, 64)

			if err != nil {
				return err, nil
			}

			return nil, &CoreRes{RefVal: reflect.ValueOf(f), Val: f}

		} else {
			//check for floats
			v, err := strconv.Atoi(strData)
			if err != nil {
				return err, nil
			}

			return nil, &CoreRes{RefVal: reflect.ValueOf(v), Val: v}
		}

	case ibStr == "l":
		var arr []any
		*i++
		for e[*i] != 'e' {
			err, val := unMarshalCore(e, val, i)
			if err != nil {
				return err, nil
			}
			arr = append(arr, val.Val)
		}
		*i++
		return nil, &CoreRes{RefVal: reflect.ValueOf(arr), Val: arr}

	case ibStr == "d":
		*i++
		newStructP := reflect.New(val.RefVal.Elem().Type())
		newStruct := newStructP.Elem()
		for e[*i] != 'e' {
			//explicit call since value and key are sequential
			err, key := unMarshalCore(e, val, i)

			if err != nil {
				return err, nil
			}

			ok, name, kind := structHasProp(newStruct.Type(), reflect.ValueOf(key.Val).String())

			if ok {
				if kind == reflect.Struct {
					newVal := newStruct.FieldByName(name)
					newValP := reflect.New(newVal.Type())
					cr := CoreRes{
						RefVal: newValP,
					}
					err, value := unMarshalCore(e, cr, i)
					if err != nil {
						return err, nil
					}
					newStruct.FieldByName(name).Set(value.RefVal)
				} else {
					err, value := unMarshalCore(e, val, i)
					if err != nil {
						return err, nil
					}
					newStruct.FieldByName(name).Set(reflect.ValueOf(value.Val))
				}
			}
		}
		*i++
		return nil, &CoreRes{RefVal: newStruct}
	default:
		return fmt.Errorf("invalid bencode text"), nil
	}
}

func structHasProp(st reflect.Type, key string) (bool, string, reflect.Kind) {
	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get("bencode")
		if tag == key {
			return true, st.Field(i).Name, st.Field(i).Type.Kind()
		}
	}
	return false, "", reflect.Struct
}
