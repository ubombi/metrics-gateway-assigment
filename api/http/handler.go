package http

import (
	"encoding/json"

	"github.com/ubombi/timeseries/api"
	"github.com/valyala/fasthttp"
)

type Service struct {
	Storage api.Storage
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

	if e.EventType == "" || e.Ts == 0 {
		ctx.Response.SetStatusCode(400)
		return
	}

	err := s.Storage.Store(api.Event{
		Type:   e.EventType,
		Ts:     e.Ts,
		Params: e.Params,
	})
	if err != nil {
		ctx.Response.SetStatusCode(500)
	}
}
