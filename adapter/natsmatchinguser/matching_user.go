package natsmatchinguser

import (
	"fmt"
	"gameapp/entity"
)

func (b Broker) Publish(event entity.Event, payload string) {
	fmt.Println("payload", payload)
	err := b.Conn.Publish(string(event), []byte(payload))
	if err != nil {
		return
	}
}
