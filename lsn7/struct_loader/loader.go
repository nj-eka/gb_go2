package struct_loader

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	TagName     = "name"
	TagRequired = "required"
	TagDefault  = "default"
)

var (
	ErrorLoader                = errors.New("reflect loader error")
	ErrorInvalidReceiver       = fmt.Errorf("invalid receiver: %w", ErrorLoader)
	ErrorInvalidProvider       = fmt.Errorf("invalid provider: %w", ErrorLoader)
	ErrorIsNil                 = fmt.Errorf("nil: %w", ErrorInvalidReceiver)
	ErrorNotSettable           = fmt.Errorf("not settable: %w", ErrorInvalidReceiver)
	ErrorInvalidDefaultSetting = fmt.Errorf("invalid default setting: %w", ErrorInvalidReceiver)
	ErrorNotConvertable        = fmt.Errorf("not convertable: %w", ErrorInvalidProvider)
	ErrorFieldRequired         = fmt.Errorf("required: %w", ErrorInvalidProvider)
)

func getFieldName(field reflect.StructField) (fieldName string) {
	fieldName = field.Name
	if tagName := field.Tag.Get(TagName); tagName != "" {
		fieldName = tagName
	}
	return
}

func isRequired(field reflect.StructField, defaultRequired bool) bool {
	if value := field.Tag.Get(TagRequired); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultRequired
}

func lookupDefaultValue(field reflect.StructField, target reflect.Value) (value interface{}, ok bool) {
	if data, ok := field.Tag.Lookup(TagDefault); ok {
		if rValue, err := parse(target, data); err == nil {
			return rValue.Interface(), true
		} else {
			panic(fmt.Errorf("[%s] %w - %v", field.Name, ErrorInvalidDefaultSetting, err))
		}
	}
	return
}

func getPrefixedMap(prefix string, values map[string]interface{}) map[string]interface{} {
	if prefix == "" {
		return values
	}
	pMap := make(map[string]interface{})
	for key, value := range values {
		if strings.HasPrefix(key, prefix) {
			pMap[strings.TrimPrefix(key, prefix)] = value
		}
	}
	return pMap
}

func parse(target reflect.Value, value interface{}) (rValue reflect.Value, err error) {
	switch v := value.(type){
	case string:
		switch target.Interface().(type) {// Kind() base types only
		case string:
			rValue = reflect.ValueOf(v)
		case time.Duration:
			var td time.Duration
			if td, err = time.ParseDuration(v); err == nil {
				rValue = reflect.ValueOf(td)
			}
		case bool:
			var b bool
			if b, err = strconv.ParseBool(v); err == nil {
				rValue = reflect.ValueOf(b)
			}
		case int, int8, uint16, int32, int64:
			var i int64
			if i, err = strconv.ParseInt(v, 0, 0); err == nil{
				rValue = reflect.ValueOf(i).Convert(target.Type())
			}
		case float32, float64:
			var f float64
			if f, err = strconv.ParseFloat(v, 0); err == nil{
				rValue = reflect.ValueOf(f).Convert(target.Type())
			}
			//todo: add cases for time, complex, ...
		default:
			switch target.Kind(){
			case reflect.Map: // todo: add recursive parsing / current workaround below
				output := reflect.MakeMap(target.Type()).Interface()
				if err = json.Unmarshal([]byte(v), &output); err == nil {
					rValue = reflect.ValueOf(output)
				}
			case reflect.Slice: // --//--
				output := reflect.MakeSlice(target.Type(), 0, len(strings.Split(v, ","))+1).Interface()
				if err = json.Unmarshal([]byte(v), &output); err == nil {
					rValue = reflect.ValueOf(output)
				}
			}
		}
	default:
		switch reflect.ValueOf(value).Kind(){
		case reflect.Slice: // []...[]interface{]
			vv := reflect.ValueOf(value)
			output := reflect.MakeSlice(target.Type(), vv.Len(), vv.Len())
			for i := 0; i < vv.Len(); i++ {
				var rElm reflect.Value
				if rElm, err = parse(output.Index(i), vv.Index(i).Interface()); err == nil{
					output.Index(i).Set(rElm)
				} else {
					panic(err)
				}
			}
			rValue = output
		// todo: add map of map ... support
		default:
			rValue = reflect.ValueOf(v).Convert(target.Type())
		}

	}
	return
}

func setValue(target reflect.Value, value interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v : %w", r, ErrorNotConvertable)
		}
	}()
	rValue := reflect.ValueOf(value)
	rValueType := rValue.Type()
	targetType := target.Type()
	if !rValueType.AssignableTo(targetType) {
		if rValueType.ConvertibleTo(targetType) {
			rValue = rValue.Convert(targetType)
		} else {
			if target.Kind() == rValue.Kind() {
				switch targetType.Kind() {
				case reflect.Map:
					output := target.Interface()
					if target.IsNil() {
						output = reflect.MakeMap(targetType).Interface()
					}
					if err = mapstructure.Decode(value, &output); err != nil { // current workaround
						panic(err)
					}
					rValue = reflect.ValueOf(output)
				case reflect.Slice:
					output := reflect.MakeSlice(targetType, rValue.Len(), rValue.Len())
					for i := 0; i < rValue.Len(); i++ {
						var rElm reflect.Value
						if rElm, err = parse(output.Index(i), rValue.Index(i).Interface()); err == nil{
							output.Index(i).Set(rElm)
						} else {
							panic(err)
						}
					}
					rValue = output
				}
			} else {
				if rValue, err = parse(target, value); err != nil{
					panic(err)
				}
			}
		}
	}
	target.Set(rValue)
	return nil
}

func Load(in interface{}, values map[string]interface{}, defaultRequired bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if rerr, ok := r.(error); ok && errors.Is(rerr, ErrorLoader) {
				err = rerr
			} else {
				err = fmt.Errorf("%v: %w", r, ErrorLoader)
			}
		}
	}()
	if in == nil {
		return ErrorIsNil
	}
	rIn := reflect.ValueOf(in)
	if rIn.Kind() == reflect.Ptr {
		rIn = rIn.Elem()
	} else {
		return fmt.Errorf("ptr to struct expected: %w", ErrorNotSettable)
	}
	if rIn.Kind() != reflect.Struct {
		return fmt.Errorf("struct expected: %w", ErrorInvalidReceiver)
	}
	for fieldIndex := 0; fieldIndex < rIn.NumField(); fieldIndex++ {
		fieldStruct := rIn.Type().Field(fieldIndex)
		fieldName := getFieldName(fieldStruct)
		fieldValue := rIn.Field(fieldIndex)
		fieldValueType := fieldValue.Type() // == rIn.Type().Field(fieldIndex).Type
		if !fieldValue.CanSet() {
			return fmt.Errorf("field [%s] of type [%s]: %w", fieldName, fieldValueType, ErrorNotSettable)
		}
		if value, ok := values[fieldName]; ok {
			if err := setValue(fieldValue, value); err != nil {
				return fmt.Errorf("[%s] set value [%v] failed: %w", fieldName, value, err)
			}
		} else {
			if fieldValueType.Kind() == reflect.Struct {
				if prefixedMap := getPrefixedMap(fieldName+".", values); len(prefixedMap) > 0 {
					if !fieldValue.IsValid() {
						fieldValue = reflect.New(fieldValueType)
					} else {
						fieldValue = reflect.NewAt(fieldValueType, unsafe.Pointer(fieldValue.UnsafeAddr()))
					}
					if err := Load(fieldValue.Interface(), prefixedMap, false); err == nil {
						continue
					}
					return err
				}
			}
			if !fieldValue.IsZero() { // IsValid()
				continue
			}
			if value, ok := lookupDefaultValue(fieldStruct, fieldValue); ok {
				if err := setValue(fieldValue, value); err != nil {
					return fmt.Errorf("[%s] set value [%v] failed: %w", fieldName, value, err)
				}
			} else {
				if isRequired(fieldStruct, defaultRequired) {
					return fmt.Errorf("missing value for required field [%s]: %w", fieldName, ErrorFieldRequired)
				}
			}
		}
	}
	return
}
