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
		fmt.Println("res", res, res.CanConvert(destVal.Elem().Type()))
		if res.CanConvert(destVal.Elem().Type().Elem()) {
			res = res.Convert(destVal.Elem().Type().Elem())
			newSlice = reflect.Append(newSlice, res)
		}
	}
	destVal.Elem().Set(newSlice)
	return nil
}

func UnMarshal(e []byte, val any) (err error, v any) {
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return fmt.Errorf("cannot pass a value instead of a reference"), nil
	}

	read := 0
	err, res := unMarshalCore(e, val, &read)

	if err != nil {
		return err, nil
	}

	valRef := reflect.ValueOf(val)
	resRef := reflect.ValueOf(res)

	if valRef.Elem().Kind() != resRef.Kind() {
		return fmt.Errorf("type mismatched"), nil
	}

	// convertSlice(resRef, val)
	// valRef.Elem().Set(resRef)

	return nil, res
}

func unMarshalCore(e []byte, val any, i *int) (err error, va any) {
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
		strVal := e[*i : *i+skip]
		// NOTE
		*i += skip
		return nil, string(strVal)

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

			return nil, f

		} else {
			//check for floats
			v, err := strconv.Atoi(strData)
			if err != nil {
				return err, nil
			}

			return nil, v
		}
	case ibStr == "l":
		var arr []any
		*i++
		for e[*i] != 'e' {
			err, val := unMarshalCore(e, val, i)
			if err != nil {
				return err, nil
			}
			arr = append(arr, val)
		}
		*i++
		return nil, arr
	case ibStr == "d":
		*i++
		newStructP := reflect.New(reflect.TypeOf(val).Elem())
		newStruct := newStructP.Elem()
		for e[*i] != 'e' {
			//explicit call since value and key are sequential
			err, key := unMarshalCore(e, val, i)
			err, val := unMarshalCore(e, val, i)
			if err != nil {
				fmt.Println(err)
				return err, nil
			}
			if ok, name := doesStructHasProp(newStruct.Type(), reflect.ValueOf(key).String()); ok {
				newStruct.FieldByName(name).Set(reflect.ValueOf(val))
			}
		}

		//build the struct
		*i++
		return nil, newStruct 

		//handle dict
	default:
		return fmt.Errorf("unkown initial byte, exiting"), nil
	}
}

func doesStructHasProp(st reflect.Type, key string) (bool, string) {
	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag.Get("bencode")
		if tag == key {
			return true, st.Field(i).Name
		}
	}
	return false, ""
}
