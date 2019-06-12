package api

import (
	context "context"
	"io"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/ubombi/timeseries/storage"
)

var emptyResp = empty.Empty{}

type Server struct {
	Storage storage.Interface
}

func (s *Server) StoreEvent(ctx context.Context, e *Event) (*empty.Empty, error) {
	err := s.Storage.Store(storage.Event{
		Type:   e.EventType,
		Ts:     e.Ts,
		Params: MapFromProtoStruct(e.Params),
	})
	return &emptyResp, err
}
func (s *Server) StreamEvents(stream EventService_StreamEventsServer) error {
	for {
		e, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = s.Storage.Store(storage.Event{
			Type:   e.EventType,
			Ts:     e.Ts,
			Params: MapFromProtoStruct(e.Params),
		})
		if err != nil {
			return err
		}

	}
}
