package clickhouse

import (
	uuid "github.com/satori/go.uuid"
	"github.com/ubombi/timeseries/api"
)

// event is internal representation of api.Event used to Insert data into clickhouse
//
// Untyped params grouped by their type and saved as `Nested` type.
type event struct {
	EventType         string
	Ts                int64
	StringParamNames  []string
	StringParamValues []string
	IntParamNames     []string
	IntParamValues    []int64
	FloatParamNames   []string
	FloatParamValues  []float64

	IP  string
	UID uuid.UUID
}

// convertToInternal converts Event to internal representation that can be used in clickhouse queries. Depends on actual table schema.
func convertToInternal(se api.Event) (e event) {
	e.EventType = se.Type
	e.Ts = se.Ts

	for name, untyped := range se.Params {
		// add predefined params here. Example below
		if name == "uid" {
			if v, ok := untyped.(string); ok {
				e.UID = uuid.FromStringOrNil(v)
				continue
			}
		}

		// due to clickhouse limitations each param value type should be saved separately
		switch value := untyped.(type) {
		case string:
			e.StringParamValues = append(e.StringParamValues, value)
			e.StringParamNames = append(e.StringParamNames, name)
		case int, int32, int64:
			e.IntParamValues = append(e.IntParamValues, toInt64(value))
			e.IntParamNames = append(e.IntParamNames, name)
		case float32, float64:
			e.FloatParamValues = append(e.FloatParamValues, toFloat64(value))
			e.FloatParamNames = append(e.FloatParamNames, name)
		}
	}
	return
}

// toInt64 helper function
func toInt64(number interface{}) int64 {
	switch v := number.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	default:
		return 0
	}
}

// toFloat64 helper function
func toFloat64(number interface{}) float64 {
	switch v := number.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	default:
		return 0
	}
}
