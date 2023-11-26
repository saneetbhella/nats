package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

const subject = "foo"

type natsServer struct {
	connection *nats.Conn
}

func main() {
	fmt.Println("Opening...")

	nc := openConnection().connection

	// Although there are no subscribers this will be published successfully
	err := nc.Publish(subject, []byte("Hello World"))
	if err != nil {
		panic(err)
	}

	// Simple Sync Subscriber
	sub, _ := nc.SubscribeSync(subject)

	// For a synchronous subscription we need to fetch the next message
	// Publish occurred before the subscription was established so will timeout
	msg, _ := sub.NextMsg(10 * time.Millisecond)

	err = nc.Publish(subject, []byte("Hello World"))
	if err != nil {
		panic(err)
	}

	msg, _ = sub.NextMsg(10 * time.Millisecond)
	if msg != nil {
		fmt.Printf(
			"--Received a message--\n"+
				"Subject: %s\n"+
				"Message: %s\n",
			msg.Subject,
			string(msg.Data),
		)
	}

	fmt.Println("Closing...")
}

func openConnection() *natsServer {
	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	return &natsServer{
		connection: nc,
	}
}

func (nc *natsServer) closeConnection() {
	// Close() not needed if this is called.
	// Drain is a safe way to ensure all buffered messages that were published are sent
	// and all buffered messages received on a subscription are processed
	err := nc.connection.Drain()
	if err != nil {
		panic(err)
	}
}
