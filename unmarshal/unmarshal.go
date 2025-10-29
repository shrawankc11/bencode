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

func UnMarshal(e []byte, val any) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return fmt.Errorf("function expects a apointer received value")
	}

	read := 0
	valRef := reflect.ValueOf(val).Elem()
	corRes := CoreRes{
		RefVal: valRef,
	}
	err, res := unMarshalCore(e, corRes, &read)

	if err != nil {
		return err
	}

	if valRef.Kind() != res.RefVal.Kind() {
		return fmt.Errorf("type mismatched")
	}

	valRef.Set(res.RefVal)

	return nil
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

		*i++
		strVal := string(e[*i : *i+skip])
		*i += skip

		return nil, &CoreRes{RefVal: reflect.ValueOf(strVal)}

	case ibStr == "i":
		*i++
		strData := ""
		for ; e[*i] != 'e'; *i++ {
			strData += string(e[*i])
		}
		*i++

		//fix this float check
		if i := strings.IndexByte(strData, '.'); i != -1 {
			f, err := strconv.ParseFloat(strData, 64)

			if err != nil {
				return err, nil
			}

			return nil, &CoreRes{RefVal: reflect.ValueOf(f)}

		} else {
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
			cr := CoreRes{
				RefVal: val.RefVal,
			}
			if val.RefVal.Type().Elem().Kind() == reflect.Slice {
				cr.RefVal = reflect.MakeSlice(val.RefVal.Type().Elem(), val.RefVal.Len(), val.RefVal.Cap())
			} else if val.RefVal.Type().Elem().Kind() == reflect.Struct {
				nestedStruct := reflect.New(val.RefVal.Type().Elem())
				cr.RefVal = nestedStruct.Elem()
			}
			err, val := unMarshalCore(e, cr, i)
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
					cr := CoreRes{
						RefVal: val.RefVal,
					}
					if kind == reflect.Slice {
						newVal := newStruct.FieldByName(name)
						newValP := reflect.New(newVal.Type())
						cr.RefVal = newValP.Elem()
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
