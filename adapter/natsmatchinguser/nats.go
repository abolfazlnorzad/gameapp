package natsmatchinguser

import "github.com/nats-io/nats.go"

type Broker struct {
	Conn *nats.Conn
}

func New(address string) Broker {
	c, err := nats.Connect(address)
	if err != nil {
		panic(err)
	}
	return Broker{Conn: c}
}
