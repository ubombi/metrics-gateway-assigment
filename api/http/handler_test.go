package http

import (
	"testing"

	"github.com/ubombi/timeseries/storage"
	"github.com/valyala/fasthttp"
)

func TestService_HandlerFunc(t *testing.T) {
	type fields struct {
		Storage storage.Interface
	}
	type args struct {
		ctx *fasthttp.RequestCtx
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: tt.fields.Storage,
			}
			s.HandlerFunc(tt.args.ctx)
		})
	}
}
