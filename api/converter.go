package api

import (
	pstruct "github.com/golang/protobuf/ptypes/struct"
)

func MapFromProtoStruct(pr *pstruct.Struct) map[string]interface{} {
	ret := make(map[string]interface{}, len(pr.Fields))
	for k, v := range pr.Fields {
		ret[k] = FieldFromProtoValue(v)
	}
	return ret

}

func ProtoStructFromMap(mp map[string]interface{}) *pstruct.Struct {
	ret := &pstruct.Struct{
		Fields: make(map[string]*pstruct.Value, len(mp)),
	}
	for k, v := range mp {
		ret.Fields[k] = ProtoValueFromField(v)
	}
	return ret

}

// fieldFromProtoValue converts `google.protobuf.Struct` value into untyped interface{} value.  protobuf Struct is used to send untyped params values within typed protobuf protocol. Returns nil for unsupported types
func FieldFromProtoValue(v *pstruct.Value) interface{} {
	switch typed := v.GetKind().(type) {
	case *pstruct.Value_NullValue:
		return typed.NullValue
	case *pstruct.Value_BoolValue:
		return typed.BoolValue
	case *pstruct.Value_NumberValue:
		return typed.NumberValue
	case *pstruct.Value_StringValue:
		return typed.StringValue
	//case *pstruct.Value_StructValue:
	//return MapFromProtoStruct(typed.Value_StructValue)
	//case *pstruct.Value_ListValue:
	//return
	default:
		return nil
	}
}

// ProtoValueFromField converts value into protobuf.Value struct. Returns null_value for unsupported types
func ProtoValueFromField(f interface{}) *pstruct.Value {
	value := &pstruct.Value{}

	if f == nil {
		value.Kind = &pstruct.Value_NullValue{}
		return value
	}

	switch typed := f.(type) {
	case int:
		value.Kind = &pstruct.Value_NumberValue{
			NumberValue: float64(typed),
		}
	case float32:
		value.Kind = &pstruct.Value_NumberValue{
			NumberValue: float64(typed),
		}
	case float64:
		value.Kind = &pstruct.Value_NumberValue{
			NumberValue: typed,
		}
	case int32:
		value.Kind = &pstruct.Value_NumberValue{
			NumberValue: float64(typed),
		}
	case int64:
		value.Kind = &pstruct.Value_NumberValue{
			NumberValue: float64(typed),
		}
	case bool:
		value.Kind = &pstruct.Value_BoolValue{
			BoolValue: typed,
		}
	case string:
		value.Kind = &pstruct.Value_StringValue{
			StringValue: typed,
		}

	}
	return value
}
