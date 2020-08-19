package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

type payload struct {
	Name      string `json:"name"`
	URLPrefix string `json:"url_prefix"`
}

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
	namespace := "container"
	if val, ok := os.LookupEnv("namespace"); ok {
		namespace = val
	}
	if err := createService(payload.Name, namespace); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Cannot create Service: %s", err)
		log.Println(errMsg)
		w.Write([]byte(errMsg))
		return
	}

	subject := "nats-create-virtual-service"
	if val, ok := os.LookupEnv("topic_create_virtual_service"); ok {
		subject = val
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not marshal json", http.StatusInternalServerError)
		return
	}

	if err := sendViaNats(subject, payloadBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg := fmt.Sprintf("Can not connect/publish to message queue: %s", err)
		w.Write([]byte(errMsg))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendViaNats(subject string, msg []byte) error {
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
