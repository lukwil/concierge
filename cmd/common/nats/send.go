package nats

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

// Send sends given message to nats subject.
func Send(subject string, msg []byte) error {
	natsURL := nats.DefaultURL
	if val, ok := os.LookupEnv("nats_url"); ok {
		natsURL = val
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer nc.Close()

	if err := nc.Publish(subject, msg); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Published %d bytes to: %q\n", len(msg), subject)
	return nil
}
