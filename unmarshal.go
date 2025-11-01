package bencode 

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type coreRes struct {
	RefVal reflect.Value
}

func UnMarshal(e []byte, val any) (err error) {
	if reflect.TypeOf(val).Kind() != reflect.Pointer {
		return fmt.Errorf("function expects a pointer received value")
	}
	read := 0
	valRef := reflect.ValueOf(val).Elem()
	corRes := coreRes{
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

func unMarshalCore(e []byte, val coreRes, i *int) (error, *coreRes) {
	initialByte := e[*i: *i+1]
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

		return nil, &coreRes{RefVal: reflect.ValueOf(strVal)}

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

			return nil, &coreRes{RefVal: reflect.ValueOf(f)}

		} else {
			v, err := strconv.Atoi(strData)
			if err != nil {
				return err, nil
			}

			return nil, &coreRes{RefVal: reflect.ValueOf(v)}
		}

	case ibStr == "l":
		var arr = reflect.MakeSlice(val.RefVal.Type(), val.RefVal.Len(), val.RefVal.Cap())
		*i++
		for e[*i] != 'e' {
			cr := coreRes{
				RefVal: val.RefVal,
			}
			/**
			* Array elements can be of type slice or struct.
			* cr.RefVal should point to the element's type.
			* For normal data types no need to update cr.RefVal since it is not used.
			 */

			currentElement := val.RefVal.Type().Elem()

			if currentElement.Kind() == reflect.Slice {
				cr.RefVal = reflect.MakeSlice(currentElement, val.RefVal.Len(), val.RefVal.Cap())
			} else if currentElement.Kind() == reflect.Struct {
				nestedStruct := reflect.New(currentElement)
				cr.RefVal = nestedStruct.Elem()
			}
			err, val := unMarshalCore(e, cr, i)
			if err != nil {
				return err, nil
			}
			arr = reflect.Append(arr, val.RefVal)
		}
		*i++
		return nil, &coreRes{RefVal: arr}

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
				cr := coreRes{
					RefVal: val.RefVal,
				}

				/*
				* Struct element's data type can be slice or struct.
				* In case of iterators cr.RefVal referes to the elements data type.
				* For normal data types no need to update cr.RefVal.
				 */
				switch kind {
				case reflect.Struct:
					newVal := newStruct.FieldByName(name)
					newStruct := reflect.New(newVal.Type())
					cr.RefVal = newStruct.Elem()
				case reflect.Slice:
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
		*i++
		return nil, &coreRes{RefVal: newStruct}
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
