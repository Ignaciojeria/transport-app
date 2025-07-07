package natsconn

import (
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	ioc.Registry(NewJetStream, NewConn)
}
func NewJetStream(conn *nats.Conn) (jetstream.JetStream, error) {
	return jetstream.New(conn)
}
