package nats

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/nats-io/nats-server/v2/server"
)

func startEmbeddedServer() (*server.Server, error) {
	opts := &server.Options{}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup embedded nats server")
	}

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, errors.New("nats server not ready for connections in 5 seconds")
	}

	return ns, nil
}
