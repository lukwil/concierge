package create

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		input = body
	}
	log.Println(string(input))

	var payload payload
	if err := json.Unmarshal(input, &payload); err != nil {
		log.Println(err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	payload.createStatefulSet()

	subject := "nats-create-service"
	if val, ok := os.LookupEnv("topic_create_service"); ok {
		subject = val
	}

	if err := sendViaNats(subject, string(input)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Can not connect/publish to message queue: %s", err)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello world, input was lukwil: %s", string(input))))
}

func sendViaNats(subject, msg string) error {
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

	if err := nc.Publish(subject, []byte(msg)); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Published %d bytes to: %q\n", len(msg), subject)
	return nil
}
