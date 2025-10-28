package unmarshal

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type CoreRes struct {
	RefVal reflect.Value
}

func UnMarshal(e []byte, val any) (err error, v any) {
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return fmt.Errorf("function expects a apointer received value"), nil
	}

	read := 0
	valRef := reflect.ValueOf(val).Elem()
	corRes := CoreRes{
		RefVal: valRef,
	}
	err, res := unMarshalCore(e, corRes, &read)

	if err != nil {
		return err, nil
	}

	fmt.Println("res", res)

	if res.RefVal.Kind() == reflect.Struct {
		for i := 0; i < res.RefVal.Type().NumField(); i++ {
			field := res.RefVal.Type().Field(i).Name
			fmt.Println("field", field)
		}
	}

	if valRef.Kind() != res.RefVal.Kind() {
		return fmt.Errorf("type mismatched"), nil
	}

	// convertSlice(resRef, val)
	valRef.Set(res.RefVal)

	return nil, res.RefVal
}

func unMarshalCore(e []byte, val CoreRes, i *int) (error, *CoreRes) {
	fmt.Println("valRef", val.RefVal, val.RefVal.Type())
	var err error

	reader := bytes.NewReader(e[*i:])
	fmt.Println(string(e[*i:]))
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

		return nil, &CoreRes{RefVal: reflect.ValueOf(strVal)}

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

			return nil, &CoreRes{RefVal: reflect.ValueOf(f)}

		} else {
			//check for floats
			v, err := strconv.Atoi(strData)
			if err != nil {
				return err, nil
			}

			return nil, &CoreRes{RefVal: reflect.ValueOf(v)}
		}

	case ibStr == "l":
		var arr = reflect.MakeSlice(val.RefVal.Type(), val.RefVal.Len(), val.RefVal.Cap())
		*i++
		for e[*i] != 'e' {
			err, val := unMarshalCore(e, val, i)
			if err != nil {
				return err, nil
			}
			arr = reflect.Append(arr, val.RefVal)
		}
		*i++
		return nil, &CoreRes{RefVal: arr}

	case ibStr == "d":
		*i++
		newStructP := reflect.New(val.RefVal.Type())
		newStruct := newStructP.Elem()
		for e[*i] != 'e' {
			//explicit call since value and key are sequential
			err, key := unMarshalCore(e, val, i)

			if err != nil {
				return err, nil
			}

			ok, name, kind := structHasProp(newStruct.Type(), key.RefVal.String())

			if ok {
				if kind == reflect.Struct {
					newVal := newStruct.FieldByName(name)
					newValP := reflect.New(newVal.Type())
					cr := CoreRes{
						RefVal: newValP.Elem(),
					}
					err, value := unMarshalCore(e, cr, i)
					if err != nil {
						return err, nil
					}
					newStruct.FieldByName(name).Set(value.RefVal)
				} else {
					if kind == reflect.Slice {
						fmt.Println("got slice", name, kind, newStruct)
					}
					newVal := newStruct.FieldByName(name)
					newValP := reflect.New(newVal.Type())
					cr := CoreRes{
						RefVal: newValP.Elem(),
					}
					err, value := unMarshalCore(e, cr, i)
					if err != nil {
						return err, nil
					}
					newStruct.FieldByName(name).Set(value.RefVal)
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
