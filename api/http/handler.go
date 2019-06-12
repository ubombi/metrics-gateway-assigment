package http

import (
	"encoding/json"

	"github.com/ubombi/timeseries/storage"
	"github.com/valyala/fasthttp"
)

type Service struct {
	Storage storage.Interface
}

func (s *Service) HandlerFunc(ctx *fasthttp.RequestCtx) {
	var e struct {
		EventType string                 `json:"event_type"`
		Ts        int64                  `json:"ts"`
		Params    map[string]interface{} `json:"params"`
	}
	if err := json.Unmarshal(ctx.PostBody(), &e); err != nil {
		ctx.Response.SetStatusCode(400)
		return
	}

	// TODO: validation here
	err := s.Storage.Store(storage.Event{
		Type:   e.EventType,
		Ts:     e.Ts,
		Params: e.Params,
	})
	if err != nil {
		ctx.Response.SetStatusCode(500)
	}
}
