package grpc

import (
	"reflect"
	"testing"

	p "github.com/golang/protobuf/ptypes/struct"
	pstruct "github.com/golang/protobuf/ptypes/struct"
)

var struct1 = &p.Struct{
	Fields: map[string]*p.Value{
		"field1": &p.Value{
			Kind: &p.Value_StringValue{
				StringValue: "string1",
			},
		},
		"field2": &p.Value{
			Kind: &p.Value_StringValue{
				StringValue: "string2",
			},
		},
	},
}

var struct2 = &p.Struct{
	Fields: map[string]*p.Value{
		"field1": &p.Value{
			Kind: &p.Value_NumberValue{
				NumberValue: 99999.99999,
			},
		},
		"field2": &p.Value{
			Kind: &p.Value_NumberValue{
				NumberValue: 0,
			},
		},
	},
}

var struct3 = &p.Struct{
	Fields: map[string]*p.Value{
		"str": &p.Value{
			Kind: &p.Value_StringValue{
				StringValue: "some string 1",
			},
		},
		"float": &p.Value{
			Kind: &p.Value_NumberValue{
				NumberValue: 12.49,
			},
		},
		"bool": &p.Value{
			Kind: &p.Value_BoolValue{
				BoolValue: true,
			},
		},
	},
}

var map1 = map[string]interface{}{
	"field1": 99999.99999,
	"field2": 0.0,
}

var map2 = map[string]interface{}{
	"field1": "string1",
	"field2": "string2",
}

var map3 = map[string]interface{}{
	"str":   "some string 1",
	"float": 12.49,
	"bool":  true,
}

func TestMapFromProtoStruct(t *testing.T) {
	tests := []struct {
		name string
		pr   *p.Struct
		want map[string]interface{}
	}{
		{
			name: "float map",
			pr:   struct2,
			want: map1,
		},
		{
			name: "string map",
			pr:   struct1,
			want: map2,
		},
		{
			name: "interface{} map",
			pr:   struct3,
			want: map3,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapFromProtoStruct(tt.pr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapFromProtoStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtoStructFromMap(t *testing.T) {
	tests := []struct {
		name string
		mp   map[string]interface{}
		want *pstruct.Struct
	}{
		{
			name: "float map",
			mp:   map1,
			want: struct2,
		},
		{
			name: "string map",
			mp:   map2,
			want: struct1,
		},
		{
			name: "interface{} map",
			mp:   map3,
			want: struct3,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProtoStructFromMap(tt.mp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoStructFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
