package grpc

import (
	context "context"
	"io"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/ubombi/timeseries/storage"
)

var emptyResp = &empty.Empty{}

type Server struct {
	Storage storage.Interface
}

func (s *Server) StoreEvent(ctx context.Context, e *Event) (*empty.Empty, error) {
	err := s.Storage.Store(storage.Event{
		Type:   e.EventType,
		Ts:     e.Ts,
		Params: MapFromProtoStruct(e.Params),
	})
	return emptyResp, err
}

// StreamEvents stores events from a stream
func (s *Server) StreamEvents(stream EventService_StreamEventsServer) error {
	// Close stream in case of error or connection lost
	// Opened stream blocks shutdown process
	defer stream.SendAndClose(emptyResp)

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
