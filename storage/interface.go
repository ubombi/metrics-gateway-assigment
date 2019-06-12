package storage

type UnixTimestamp = int64

type Event struct {
	Type   string
	Ts     UnixTimestamp
	Params map[string]interface{}
}

type Interface interface {
	Store(Event) error
}
