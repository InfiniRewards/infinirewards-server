package tests

import (
	"fmt"
	"infinirewards/nats"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

func setupNATSServer() error {

	// Test NATS connection
	if nats.NC != nil {
		// Test NATS connection is alive
		if nats.NC.IsConnected() {
			return nil
		}
	}

	opts := &natsserver.Options{
		Host:      "127.0.0.1",
		Port:      4222,
		NoLog:     false,
		Debug:     true,
		Trace:     true,
		NoSigs:    true,
		JetStream: true,
	}

	var err error
	natsServer, err = natsserver.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create NATS server: %v", err)
	}

	go natsServer.Start()

	if !natsServer.ReadyForConnections(4 * time.Second) {
		return fmt.Errorf("NATS server failed to start")
	}

	// Initialize NATS module
	if err := nats.ConnectNatsTest(); err != nil {
		return fmt.Errorf("failed to initialize NATS module: %v", err)
	}

	return nil
}
