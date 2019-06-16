package api

type UnixTimestamp = int64

type Event struct {
	Type   string
	Ts     UnixTimestamp
	Params map[string]interface{}
}

type Storage interface {
	Store(Event) error
}
